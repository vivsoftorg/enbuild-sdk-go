package enbuild_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/vivsoftorg/enbuild-sdk-go/pkg/enbuild"
	// "github.com/vivsoftorg/enbuild-sdk-go/internal/request" // Not directly used, but client uses it
)

const (
	// This is the path the client will prefix to stack routes.
	// The mock server needs to expect this full path.
	expectedApiVersionPath = "/enbuild-bk/api/v1/"
)

func TestListStacks(t *testing.T) {
	type testCase struct {
		name              string
		page              int
		limit             int
		searchTerm        string
		mockServerHandler func(w http.ResponseWriter, r *http.Request)
		expectedStacks    []*enbuild.Stack
		expectError       bool
		expectedErrorMsg  string
	}

	testCases := []testCase{
		{
			name:       "BasicFetch - No search term",
			page:       0,
			limit:      10,
			searchTerm: "",
			mockServerHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET, got %s", r.Method)
					http.Error(w, "bad method", http.StatusMethodNotAllowed)
					return
				}
				expectedPath := expectedApiVersionPath + "stacks?page=0&limit=10&search="
				if r.URL.String() != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.String())
					http.Error(w, "bad path", http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string][]*enbuild.Stack{
					"data": {{ID: "1", Name: "Stack1"}, {ID: "2", Name: "Stack2"}},
				})
			},
			expectedStacks: []*enbuild.Stack{{ID: "1", Name: "Stack1"}, {ID: "2", Name: "Stack2"}},
			expectError:    false,
		},
		{
			name:       "WithSearchTerm",
			page:       1,
			limit:      5,
			searchTerm: "searchMe",
			mockServerHandler: func(w http.ResponseWriter, r *http.Request) {
				expectedPath := expectedApiVersionPath + "stacks?page=1&limit=5&search=searchMe"
				if r.URL.String() != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.String())
					http.Error(w, "bad path", http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string][]*enbuild.Stack{
					"data": {{ID: "3", Name: "SearchStack"}},
				})
			},
			expectedStacks: []*enbuild.Stack{{ID: "3", Name: "SearchStack"}},
			expectError:    false,
		},
		{
			name:       "SearchTermWithSpecialChars",
			page:       0,
			limit:      10,
			searchTerm: "test&name=val", // Will be escaped to test%26name%3Dval
			mockServerHandler: func(w http.ResponseWriter, r *http.Request) {
				expectedQuery := "page=0&limit=10&search=" + url.QueryEscape("test&name=val")
				expectedPathSuffix := "stacks?" + expectedQuery
				if !strings.HasSuffix(r.URL.String(), expectedPathSuffix) {
					t.Errorf("Expected URL to end with %s, got %s", expectedPathSuffix, r.URL.String())
					http.Error(w, "bad path or query", http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string][]*enbuild.Stack{
					"data": {{ID: "4", Name: "SpecialCharStack"}},
				})
			},
			expectedStacks: []*enbuild.Stack{{ID: "4", Name: "SpecialCharStack"}},
			expectError:    false,
		},
		{
			name:       "EmptyResponseFromServer",
			page:       0,
			limit:      10,
			searchTerm: "",
			mockServerHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string][]*enbuild.Stack{"data": {}}) // Empty slice
			},
			expectedStacks: []*enbuild.Stack{}, // Expect an empty slice, not nil
			expectError:    false,
		},
		{
			name:       "ServerError_500",
			page:       0,
			limit:      10,
			searchTerm: "",
			mockServerHandler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "internal server error", http.StatusInternalServerError)
			},
			expectedStacks:   nil,
			expectError:      true,
			expectedErrorMsg: "API error: 500 Internal Server Error",
		},
		{
			name:       "MalformedJSONResponse",
			page:       0,
			limit:      10,
			searchTerm: "",
			mockServerHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, `{"data": [malformed_json]}`)
			},
			expectedStacks:   nil, // or empty slice depending on how error handling is done
			expectError:      true,
			expectedErrorMsg: "invalid character 'm' looking for beginning of value", // error from json.Decode
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tc.mockServerHandler))
			defer server.Close()

			// Must use context.Background() for NewClient as per its new signature
			ctx := context.Background()

			// Client setup - ensure BaseURL is correctly set for the test server
			// The enbuild.NewClient will append the apiVersionPath internally if not present.
			// So we give it server.URL which is like "http://127.0.0.1:PORT"
			// It will become "http://127.0.0.1:PORT/enbuild-bk/api/v1/"
			client, err := enbuild.NewClient(ctx, enbuild.WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			stacks, err := client.Stacks.ListStacks(ctx, tc.page, tc.limit, tc.searchTerm)

			if tc.expectError {
				if err == nil {
					t.Fatalf("Expected an error, but got none")
				}
				if tc.expectedErrorMsg != "" && !strings.Contains(err.Error(), tc.expectedErrorMsg) {
					t.Errorf("Expected error message to contain '%s', got '%s'", tc.expectedErrorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("Did not expect an error, but got: %v", err)
				}
				if !reflect.DeepEqual(stacks, tc.expectedStacks) {
					t.Errorf("Expected stacks %+v, got %+v", tc.expectedStacks, stacks)
				}
			}
		})
	}
}
