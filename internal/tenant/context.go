package tenant

import (
	"context"
	"errors"
)

var ErrNoTenant = errors.New("tenant ID missing from context")

type contextKey struct{}

func WithID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, contextKey{}, id)
}

func FromContext(ctx context.Context) (int64, error) {
	id, ok := ctx.Value(contextKey{}).(int64)
	if !ok || id <= 0 {
		return 0, ErrNoTenant
	}
	return id, nil
}
