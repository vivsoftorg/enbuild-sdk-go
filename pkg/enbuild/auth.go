package enbuild

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// AdminSettingsResponse represents the response from the admin settings API
type AdminSettingsResponse struct {
	Data map[string]AdminSettingData `json:"data"`
}

// AdminSettingData represents the admin settings data
type AdminSettingData struct {
	AuthMechanism string `json:"authMechanism"`
	AdminConfigs  struct {
		Keycloak struct {
			KeycloakBackendURL string `json:"KEYCLOAK_BACKEND_URL"`
			KeycloakClientID   string `json:"KEYCLOAK_CLIENT_ID"`
			KeycloakRealm      string `json:"KEYCLOAK_REALM"`
		} `json:"keycloak"`
	} `json:"adminConfigs"`
}

// KeycloakConfig holds the Keycloak configuration from admin settings
type KeycloakConfig struct {
	BackendURL string
	ClientID   string
	Realm      string
}

// KeycloakTokenResponse represents the token response from Keycloak
type KeycloakTokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

// AuthManager handles authentication and token refresh
type AuthManager struct {
	username       string
	password       string
	keycloakConfig KeycloakConfig
	accessToken    string
	refreshToken   string
	expiresAt      time.Time
	mutex          sync.RWMutex
	debug          bool
	baseURL        string
	authMechanism  string
}

// NewAuthManager creates a new AuthManager
func NewAuthManager(username, password string, debug bool, baseURL string) *AuthManager {
	// Check if debug is enabled via environment variable
	if os.Getenv("ENBUILD_DEBUG") == "true" {
		debug = true
	}

	return &AuthManager{
		username: username,
		password: password,
		debug:    debug,
		baseURL:  baseURL,
	}
}

// Initialize fetches the authentication configuration and initial token
func (am *AuthManager) Initialize() error {
	if am.debug {
		fmt.Println("DEBUG: Authenticating with username:", am.username)
	}

	// Fetch configuration from admin settings API
	if err := am.fetchAdminSettings(); err != nil {
		return fmt.Errorf("failed to fetch authentication configuration: %v", err)
	}

	// For local auth mechanism, no need to fetch a token
	if am.authMechanism == "local" {
		if am.debug {
			fmt.Printf("DEBUG: Using local authentication mechanism, no token fetch needed\n")
		}
		return nil
	}

	// Only proceed with Keycloak authentication if authMechanism is set to keycloak
	if am.authMechanism != "keycloak" {
		return fmt.Errorf("unsupported authentication mechanism: %s", am.authMechanism)
	}

	// Fetch initial token
	if err := am.fetchNewToken(); err != nil {
		return fmt.Errorf("failed to authenticate with Keycloak: %v", err)
	}

	return nil
}

// fetchAdminSettings retrieves the authentication configuration from the admin settings API
func (am *AuthManager) fetchAdminSettings() error {
	// Construct the admin settings API URL
	baseURL := am.baseURL

	// Extract the domain from the base URL
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("failed to parse base URL: %v", err)
	}

	// For development environment, use enbuild-dev.vivplatform.io
	host := parsedURL.Host
	if strings.Contains(host, "localhost") || strings.Contains(host, "127.0.0.1") {
		host = "enbuild-dev.vivplatform.io"
	}

	// Construct the admin settings URL using the appropriate domain
	// According to the requirements, the path should be /enbuild-user/api/v1/adminSettings
	adminSettingsURL := fmt.Sprintf("%s://%s/enbuild-user/api/v1/adminSettings",
		parsedURL.Scheme, host)

	if am.debug {
		fmt.Printf("DEBUG: Fetching auth config from: %s\n", adminSettingsURL)
	}

	// Make the request to the admin settings API
	resp, err := http.Get(adminSettingsURL)
	if err != nil {
		return fmt.Errorf("Failed to fetch authMechanism from ENBUILD. Please check ENBUILD_BASE_URL or network connectivity: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read admin settings response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to fetch authMechanism from ENBUILD. API returned status code %d: %s",
			resp.StatusCode, string(bodyBytes))
	}

	// Parse the response
	var adminSettings AdminSettingsResponse
	if err := json.Unmarshal(bodyBytes, &adminSettings); err != nil {
		return fmt.Errorf("failed to parse admin settings: %v", err)
	}

	// Extract the authentication configuration
	configFound := false
	for _, setting := range adminSettings.Data {
		am.authMechanism = setting.AuthMechanism

		if am.authMechanism == "keycloak" && setting.AdminConfigs.Keycloak.KeycloakBackendURL != "" {
			am.keycloakConfig.BackendURL = setting.AdminConfigs.Keycloak.KeycloakBackendURL
			am.keycloakConfig.ClientID = setting.AdminConfigs.Keycloak.KeycloakClientID
			am.keycloakConfig.Realm = setting.AdminConfigs.Keycloak.KeycloakRealm
			configFound = true
			break
		} else if am.authMechanism == "local" {
			configFound = true
			break
		}
	}

	// If no valid config was found in the response, return an error
	if !configFound {
		return fmt.Errorf("no valid authentication configuration found in admin settings")
	}

	if am.debug {
		fmt.Printf("DEBUG: Auth mechanism: %s\n", am.authMechanism)
		if am.authMechanism == "keycloak" {
			fmt.Printf("DEBUG: Using Keycloak config: URL=%s, Realm=%s, ClientID=%s\n",
				am.keycloakConfig.BackendURL,
				am.keycloakConfig.Realm,
				am.keycloakConfig.ClientID)
		} else if am.authMechanism == "local" {
			fmt.Printf("DEBUG: Using local authentication mechanism\n")
		}
	}

	return nil
}

// fetchNewToken gets a new token using username/password
func (am *AuthManager) fetchNewToken() error {
	// Ensure the backend URL has a protocol scheme
	backendURL := am.keycloakConfig.BackendURL
	if backendURL == "" {
		return fmt.Errorf("Keycloak backend URL is not set")
	}

	// Ensure URL has protocol
	if !strings.HasPrefix(backendURL, "http://") && !strings.HasPrefix(backendURL, "https://") {
		backendURL = "https://" + backendURL
	}

	// Ensure realm is not empty
	realm := am.keycloakConfig.Realm
	if realm == "" {
		return fmt.Errorf("Keycloak realm is not set")
	}

	// Ensure client ID is not empty
	clientID := am.keycloakConfig.ClientID
	if clientID == "" {
		return fmt.Errorf("Keycloak client ID is not set")
	}

	// Construct token URL according to the requirements
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token",
		backendURL,
		realm)

	if am.debug {
		fmt.Printf("DEBUG: Requesting new token from: %s\n", tokenURL)
		fmt.Printf("DEBUG: Using username: %s\n", am.username)
		fmt.Printf("DEBUG: Using client ID: %s\n", clientID)
	}

	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", clientID)
	data.Set("username", am.username)
	data.Set("password", am.password)

	return am.requestToken(tokenURL, data)
}

// refreshExpiredToken refreshes the token using the refresh token
func (am *AuthManager) refreshExpiredToken() error {
	// Ensure the backend URL has a protocol scheme
	backendURL := am.keycloakConfig.BackendURL
	if backendURL == "" {
		return fmt.Errorf("Keycloak backend URL is not set")
	}

	// Ensure URL has protocol
	if !strings.HasPrefix(backendURL, "http://") && !strings.HasPrefix(backendURL, "https://") {
		backendURL = "https://" + backendURL
	}

	// Ensure realm is not empty
	realm := am.keycloakConfig.Realm
	if realm == "" {
		return fmt.Errorf("Keycloak realm is not set")
	}

	// Ensure client ID is not empty
	clientID := am.keycloakConfig.ClientID
	if clientID == "" {
		return fmt.Errorf("Keycloak client ID is not set")
	}

	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token",
		backendURL,
		realm)

	if am.debug {
		fmt.Printf("DEBUG: Refreshing token from: %s\n", tokenURL)
	}

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", clientID)
	data.Set("refresh_token", am.refreshToken)

	return am.requestToken(tokenURL, data)
}

// requestToken makes the actual HTTP request to get or refresh a token
func (am *AuthManager) requestToken(tokenURL string, data url.Values) error {
	if am.debug {
		fmt.Printf("DEBUG: Making POST request to: %s\n", tokenURL)
	}

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Authentication with Keycloak failed. Check credentials or Keycloak settings: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Authentication with Keycloak failed. Status code %d, response: %s",
			resp.StatusCode, string(bodyBytes))
	}

	var tokenResponse KeycloakTokenResponse
	if err := json.Unmarshal(bodyBytes, &tokenResponse); err != nil {
		return fmt.Errorf("failed to decode token response: %v", err)
	}

	am.mutex.Lock()
	am.accessToken = tokenResponse.AccessToken
	am.refreshToken = tokenResponse.RefreshToken
	am.expiresAt = time.Now().Add(time.Duration(tokenResponse.ExpiresIn-30) * time.Second) // Buffer of 30 seconds
	am.mutex.Unlock()

	if am.debug {
		fmt.Printf("DEBUG: Token obtained, expires in %d seconds\n", tokenResponse.ExpiresIn)
		if len(tokenResponse.AccessToken) > 10 {
			fmt.Printf("DEBUG: Token: %s***\n", tokenResponse.AccessToken[:10])
		} else {
			fmt.Printf("DEBUG: Token: %s***\n", tokenResponse.AccessToken)
		}
		fmt.Println("DEBUG: Authentication successful!")
	}

	return nil
}

// GetToken returns a valid token, refreshing if necessary
func (am *AuthManager) GetToken() (string, error) {
	// If auth mechanism is "local", return the hardcoded token
	if am.authMechanism == "local" {
		if am.debug {
			fmt.Println("DEBUG: Using local authentication with hardcoded token")
		}
		return "enbuild_local_admin_token", nil
	}

	// If auth mechanism is not keycloak or local, return an error
	if am.authMechanism != "keycloak" {
		return "", fmt.Errorf("unsupported authentication mechanism: %s", am.authMechanism)
	}

	am.mutex.RLock()
	isExpired := time.Now().After(am.expiresAt)
	token := am.accessToken
	am.mutex.RUnlock()

	if isExpired {
		if am.debug {
			fmt.Println("DEBUG: Token expired, refreshing...")
		}
		if err := am.refreshExpiredToken(); err != nil {
			return "", fmt.Errorf("failed to refresh token: %v", err)
		}
		am.mutex.RLock()
		token = am.accessToken
		am.mutex.RUnlock()
	}

	if token == "" {
		return "", fmt.Errorf("no valid authentication token available")
	}

	return token, nil
}
