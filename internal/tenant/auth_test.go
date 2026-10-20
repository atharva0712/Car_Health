package tenant

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"car-monitoring-system/internal/logger"
)

type fakeKeyLookup struct {
	tenantID int64
	err      error
	hash     string
}

func (f *fakeKeyLookup) GetTenantIDByAPIKeyHash(_ context.Context, hash string) (int64, error) {
	f.hash = hash
	return f.tenantID, f.err
}

func TestAPIKeyGenerationAndHash(t *testing.T) {
	first, err := GenerateAPIKey()
	if err != nil {
		t.Fatal(err)
	}
	second, err := GenerateAPIKey()
	if err != nil {
		t.Fatal(err)
	}
	if len(first) != 64 || first == second {
		t.Fatalf("expected distinct 32-byte hex keys, got lengths %d and %d", len(first), len(second))
	}
	if HashAPIKey(first) == first || len(HashAPIKey(first)) != 64 {
		t.Fatal("expected SHA-256 hex digest rather than the plaintext key")
	}
}

func TestFromContextRequiresTenant(t *testing.T) {
	if _, err := FromContext(context.Background()); !errors.Is(err, ErrNoTenant) {
		t.Fatalf("expected ErrNoTenant, got %v", err)
	}
}

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		tenantID   int64
		lookupErr  error
		statusCode int
	}{
		{name: "missing key", statusCode: http.StatusUnauthorized},
		{name: "unknown key", key: "wrong", lookupErr: sql.ErrNoRows, statusCode: http.StatusUnauthorized},
		{name: "lookup failure", key: "provided", lookupErr: errors.New("database unavailable"), statusCode: http.StatusInternalServerError},
		{name: "valid key", key: "valid", tenantID: 42, statusCode: http.StatusNoContent},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lookup := &fakeKeyLookup{tenantID: test.tenantID, err: test.lookupErr}
			log := logger.NewLogger()
			called := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				id, err := FromContext(r.Context())
				if err != nil || id != 42 {
					t.Errorf("expected tenant 42 in request context, got %d and %v", id, err)
				}
				w.WriteHeader(http.StatusNoContent)
			})
			handler := NewAuthMiddleware(lookup, log).Wrap(next)

			request := httptest.NewRequest(http.MethodGet, "/", nil)
			if test.key != "" {
				request.Header.Set("X-API-Key", test.key)
			}
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)

			if response.Code != test.statusCode {
				t.Fatalf("expected status %d, got %d", test.statusCode, response.Code)
			}
			if called != (test.statusCode == http.StatusNoContent) {
				t.Fatalf("next handler called = %v, status = %d", called, response.Code)
			}
			if test.key != "" && lookup.hash != HashAPIKey(test.key) {
				t.Fatalf("lookup received hash %q, want hash of provided API key", lookup.hash)
			}
		})
	}
}
