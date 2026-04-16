package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

type contextKey string

const contextKeyRole contextKey = "role"

func JWTAuth(tokenAuth *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	verifier := jwtauth.Verifier(tokenAuth)
	authenticator := jwtauth.Authenticator(tokenAuth)

	return func(next http.Handler) http.Handler {
		withToken := verifier(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			if role, ok := claims["role"].(string); ok && role != "" {
				ctx := context.WithValue(r.Context(), contextKeyRole, role)
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		}))
		return authenticator(withToken)
	}
}

func RoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(contextKeyRole).(string)
	return role, ok
}
