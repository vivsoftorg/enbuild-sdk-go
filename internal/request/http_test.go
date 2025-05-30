package request_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	// Assuming 'request' is the alias for the package under test.
	// The actual import path will depend on the Go module structure.
	// For this example, let's use a placeholder or a common pattern.
	// If the tests are in the same package, direct access is possible (request.Client).
	// If in a _test package, then aliasing the import is common.
	// Let's assume the module path from previous examples:
	"github.com/vivsoftorg/enbuild-sdk-go/internal/request"
)

// testContextKey is a custom type for context keys to avoid collisions.
type testContextKey string

func TestNewRequestWithContext(t *testing.T) {
	baseURL, _ := url.Parse("http://localhost/api/")
	client := &request.Client{
		BaseURL:   baseURL,
		UserAgent: "test-agent",
	}

	// 1. Basic context propagation
	t.Run("ContextPropagation", func(t *testing.T) {
		const myKey = testContextKey("key1")
		ctxValue := "value1"
		ctx := context.WithValue(context.Background(), myKey, ctxValue)

		req, err := client.NewRequest(ctx, http.MethodGet, "testpath", nil)
		if err != nil {
			t.Fatalf("NewRequest returned error: %v", err)
		}

		if reqCtxValue := req.Context().Value(myKey); reqCtxValue != ctxValue {
			t.Errorf("Request context does not contain the expected value. Got %v, want %v", reqCtxValue, ctxValue)
		}
		if req.Context() != ctx {
			t.Errorf("req.Context() is not the same instance as the passed context.")
		}
	})

	// 2. Context with TokenProvider
	t.Run("ContextWithTokenProvider", func(t *testing.T) {
		const myKey = testContextKey("key2")
		ctxValue := "value2"
		ctx := context.WithValue(context.Background(), myKey, ctxValue)

		var tokenProviderCalledWith context.Context
		clientWithTokenProvider := &request.Client{
			BaseURL:   baseURL,
			UserAgent: "test-agent",
			TokenProvider: func(c context.Context) string {
				tokenProviderCalledWith = c
				return "test-token"
			},
		}

		req, err := clientWithTokenProvider.NewRequest(ctx, http.MethodGet, "testpath2", nil)
		if err != nil {
			t.Fatalf("NewRequest returned error: %v", err)
		}

		if reqCtxValue := req.Context().Value(myKey); reqCtxValue != ctxValue {
			t.Errorf("Request context does not contain the expected value. Got %v, want %v", reqCtxValue, ctxValue)
		}
		if req.Context() != ctx {
			t.Errorf("req.Context() is not the same instance as the passed context for token provider case.")
		}

		if tokenProviderCalledWith == nil {
			t.Errorf("TokenProvider was not called")
		} else {
			if tpCtxValue := tokenProviderCalledWith.Value(myKey); tpCtxValue != ctxValue {
				t.Errorf("TokenProvider was not called with the expected context. Got value %v, want %v", tpCtxValue, ctxValue)
			}
			if tokenProviderCalledWith != ctx {
			    t.Errorf("TokenProvider context is not the same instance as the passed context.")
			}
		}
		if authHeader := req.Header.Get("Authorization"); authHeader != "Bearer test-token" {
			t.Errorf("Authorization header not set correctly. Got %s", authHeader)
		}
	})

	// 3. Context with body
	t.Run("ContextWithBody", func(t *testing.T) {
		const myKey = testContextKey("key3")
		ctxValue := "value3"
		ctx := context.WithValue(context.Background(), myKey, ctxValue)

		bodyData := map[string]string{"field": "data"}
		req, err := client.NewRequest(ctx, http.MethodPost, "testpath3", bodyData)
		if err != nil {
			t.Fatalf("NewRequest returned error: %v", err)
		}

		if reqCtxValue := req.Context().Value(myKey); reqCtxValue != ctxValue {
			t.Errorf("Request context with body does not contain the expected value. Got %v, want %v", reqCtxValue, ctxValue)
		}
		if req.Context() != ctx {
			t.Errorf("req.Context() is not the same instance as the passed context for body case.")
		}
		if req.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type header not set to application/json for body. Got %s", req.Header.Get("Content-Type"))
		}
	})
}

func TestDoWithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for a custom header that NewRequest might set based on context, if desired for deeper testing.
		// For now, just ensure the request comes through.
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	serverURL, _ := url.Parse(server.URL)
	client := &request.Client{
		BaseURL:    serverURL,
		UserAgent:  "test-agent-do",
		HTTPClient: server.Client(), // Use the test server's client
	}

	const myKey = testContextKey("keyDo")
	ctxValue := "valueDo"
	ctx := context.WithValue(context.Background(), myKey, ctxValue)

	// Create request using NewRequest, which now embeds the context
	req, err := client.NewRequest(ctx, http.MethodGet, "/get", nil)
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}

	// The 'ctx' passed to Do is for potential operations *within* Do,
	// not for modifying the request's context for the HTTP call itself,
	// as that's already set by NewRequest.
	var responseData map[string]string
	resp, err := client.Do(ctx, req, &responseData)
	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	if status, ok := responseData["status"]; !ok || status != "ok" {
		t.Errorf("Response data not decoded correctly or status not ok. Got %v", responseData)
	}

	// Verify the request sent to the server had its context (implicitly via req.WithContext)
	// This is harder to directly verify from the server side without specific context propagation checks
	// (e.g. via headers if the context contained specific serializable values for that).
	// However, we've tested NewRequest sets the context on the http.Request.
	// The http.Client.Do method is responsible for using that request's context.
	// So, this test primarily ensures Do executes correctly with a context-aware request.
}

// Example of a more complex test for TokenProvider with error
func TestNewRequest_TokenProviderError(t *testing.T) {
	// This test is not explicitly required by the prompt but shows good practice
	// For now, we'll skip implementing it fully unless asked.
	t.Skip("Skipping TokenProviderError test for now as it's not strictly part of the current subtask requirements.")
}
