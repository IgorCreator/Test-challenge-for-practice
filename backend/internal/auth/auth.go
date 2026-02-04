package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type Role string

const (
	RoleMerchant Role = "merchant"
	RoleCustomer Role = "customer"
)

type Principal struct {
	Role      Role
	UserID    int64
	StoreID   int64
	StoreSlug string
	Username  string
}

type authenticator interface {
	Authenticate(ctx context.Context, username, password string) (*Principal, error)
}

type contextKey struct{}

func Middleware(authenticator authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			if !ok || strings.TrimSpace(username) == "" {
				unauthorized(w)
				return
			}
			principal, err := authenticator.Authenticate(r.Context(), username, password)
			if err != nil {
				unauthorized(w)
				return
			}
			ctx := context.WithValue(r.Context(), contextKey{}, principal)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func FromContext(ctx context.Context) (*Principal, error) {
	val := ctx.Value(contextKey{})
	if val == nil {
		return nil, errors.New("missing principal")
	}
	principal, ok := val.(*Principal)
	if !ok {
		return nil, errors.New("invalid principal")
	}
	return principal, nil
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="nimble"`)
	http.Error(w, "unauthorized", http.StatusUnauthorized)
}
