package enbuild

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	AdminConfigs struct {
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

// Initialize fetches the Keycloak configuration and initial token
func (am *AuthManager) Initialize() error {
	if am.debug {
		fmt.Println("Authenticating with username:", am.username)
	}

	// Fetch Keycloak configuration from admin settings API
	if err := am.fetchKeycloakConfig(); err != nil {
		return fmt.Errorf("failed to fetch Keycloak configuration: %v", err)
	}
	
	// Fetch initial token
	return am.fetchNewToken()
}

// fetchKeycloakConfig retrieves the Keycloak configuration from the admin settings API
func (am *AuthManager) fetchKeycloakConfig() error {
	// Construct the admin settings API URL
	baseURL := am.baseURL
	
	// Extract the domain from the base URL
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("failed to parse base URL: %v", err)
	}
	
	// Construct the admin settings URL using the same domain
	adminSettingsURL := fmt.Sprintf("%s://%s/enbuild-user/api/v1/adminSettings", 
		parsedURL.Scheme, parsedURL.Host)
	
	if am.debug {
		fmt.Printf("DEBUG: Fetching Keycloak config from: %s\n", adminSettingsURL)
	}
	
	// Make the request to the admin settings API
	resp, err := http.Get(adminSettingsURL)
	if err != nil {
		return fmt.Errorf("failed to get admin settings: %v", err)
	}
	defer resp.Body.Close()
	
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read admin settings response: %v", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get admin settings: status code %d, response: %s", 
			resp.StatusCode, string(bodyBytes))
	}
	
	// Parse the response
	var adminSettings AdminSettingsResponse
	if err := json.Unmarshal(bodyBytes, &adminSettings); err != nil {
		return fmt.Errorf("failed to parse admin settings: %v", err)
	}
	
	// Extract the Keycloak configuration
	// The admin settings response has a map with keys like "0", "1", etc.
	// We'll take the first one we find
	for _, setting := range adminSettings.Data {
		am.keycloakConfig.BackendURL = setting.AdminConfigs.Keycloak.KeycloakBackendURL
		am.keycloakConfig.ClientID = setting.AdminConfigs.Keycloak.KeycloakClientID
		am.keycloakConfig.Realm = setting.AdminConfigs.Keycloak.KeycloakRealm
		break
	}
	
	if am.keycloakConfig.BackendURL == "" || am.keycloakConfig.ClientID == "" || am.keycloakConfig.Realm == "" {
		return fmt.Errorf("incomplete Keycloak configuration in admin settings")
	}
	
	if am.debug {
		fmt.Printf("DEBUG: Using Keycloak config: URL=%s, Realm=%s, ClientID=%s\n", 
			am.keycloakConfig.BackendURL, 
			am.keycloakConfig.Realm,
			am.keycloakConfig.ClientID)
	}
	
	return nil
}

// fetchNewToken gets a new token using username/password
func (am *AuthManager) fetchNewToken() error {
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", 
		am.keycloakConfig.BackendURL, 
		am.keycloakConfig.Realm)
	
	if am.debug {
		fmt.Printf("DEBUG: Requesting new token from: %s\n", tokenURL)
		fmt.Printf("DEBUG: Using username: %s\n", am.username)
	}
	
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", am.keycloakConfig.ClientID)
	data.Set("username", am.username)
	data.Set("password", am.password)
	
	return am.requestToken(tokenURL, data)
}

// refreshExpiredToken refreshes the token using the refresh token
func (am *AuthManager) refreshExpiredToken() error {
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", 
		am.keycloakConfig.BackendURL, 
		am.keycloakConfig.Realm)
	
	if am.debug {
		fmt.Printf("DEBUG: Refreshing token from: %s\n", tokenURL)
	}
	
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", am.keycloakConfig.ClientID)
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
		return fmt.Errorf("failed to get token from Keycloak: %v", err)
	}
	defer resp.Body.Close()
	
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get token from Keycloak: status code %d, response: %s", 
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
		fmt.Println("Authentication successful!")
	}
	
	return nil
}

// GetToken returns a valid token, refreshing if necessary
func (am *AuthManager) GetToken() (string, error) {
	am.mutex.RLock()
	isExpired := time.Now().After(am.expiresAt)
	token := am.accessToken
	am.mutex.RUnlock()
	
	if isExpired {
		if am.debug {
			fmt.Println("DEBUG: Token expired, refreshing...")
		}
		if err := am.refreshExpiredToken(); err != nil {
			return "", err
		}
		am.mutex.RLock()
		token = am.accessToken
		am.mutex.RUnlock()
	}
	
	return token, nil
}
