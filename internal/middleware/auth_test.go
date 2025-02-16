package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockValidateToken(valid bool) func(string) (int, error) {
	return func(token string) (int, error) {
		if valid && token == "valid-token" {
			return 123, nil
		}
		return 0, errors.New("invalid token")
	}
}

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		validateToken  func(string) (int, error)
		expectedStatus int
	}{
		{
			name:           "Valid token",
			authHeader:     "Bearer valid-token",
			validateToken:  mockValidateToken(true),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid token",
			authHeader:     "Bearer invalid-token",
			validateToken:  mockValidateToken(false),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Missing token",
			authHeader:     "",
			validateToken:  mockValidateToken(false),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Malformed token",
			authHeader:     "Bearer",
			validateToken:  mockValidateToken(false),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				userID := r.Context().Value("userID")
				if userID == nil {
					t.Errorf("userID not found in context")
				}
				w.WriteHeader(http.StatusOK)
			})

			middleware := AuthMiddleware(tt.validateToken)
			handler := middleware(nextHandler)

			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", tt.authHeader)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestExtractToken(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
		expected   string
	}{
		{
			name:       "Valid token",
			authHeader: "Bearer valid-token",
			expected:   "valid-token",
		},
		{
			name:       "No Authorization header",
			authHeader: "",
			expected:   "",
		},
		{
			name:       "Malformed token (no Bearer)",
			authHeader: "valid-token",
			expected:   "",
		},
		{
			name:       "Malformed token (empty Bearer)",
			authHeader: "Bearer",
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", tt.authHeader)

			token := ExtractToken(req)
			if token != tt.expected {
				t.Errorf("expected token %q, got %q", tt.expected, token)
			}
		})
	}
}
