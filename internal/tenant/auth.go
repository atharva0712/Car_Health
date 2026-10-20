package tenant

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"net/http"

	"car-monitoring-system/internal/logger"
)

type KeyLookup interface {
	GetTenantIDByAPIKeyHash(context.Context, string) (int64, error)
}

type AuthMiddleware struct {
	lookup KeyLookup
	log    *logger.Logger
}

func NewAuthMiddleware(lookup KeyLookup, log *logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{lookup: lookup, log: log}
}

func HashAPIKey(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])
}

func (m *AuthMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-Key")
		if key == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		id, err := m.lookup.GetTenantIDByAPIKeyHash(r.Context(), HashAPIKey(key))
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if err != nil {
			m.log.Errorf("Failed to resolve API key: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		ctx := WithID(r.Context(), id)
		ctx = logger.WithContext(ctx, m.log.WithContext("tenant_id", id))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
