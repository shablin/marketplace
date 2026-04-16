package middleware

import "net/http"

func RBAC(allowed map[string]struct{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := RoleFromContext(r.Context())
			if !ok {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			if _, ok := allowed[role]; !ok {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func Roles(roles ...string) map[string]struct{} {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return allowed
}
