package gateway

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	middleware "github.com/shablin/marketplace/pkg/middleware"
)

type Router struct {
	handler http.Handler
}

func NewRouter(cfg Config) (*Router, error) {
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		writeJSONError(w, http.StatusNotFound, "not found")
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
	})

	tokenAuth := jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)
	r.Use(middleware.JWTAuth(tokenAuth))

	rbac := map[string]map[string]struct{}{
		"/users":         middleware.Roles("buyer", "seller", "admin"),
		"/products":      middleware.Roles("buyer", "seller", "admin"),
		"/cart":          middleware.Roles("buyer", "admin"),
		"/orders":        middleware.Roles("buyer", "seller", "admin"),
		"/payments":      middleware.Roles("buyer", "admin"),
		"/notifications": middleware.Roles("buyer", "seller", "admin"),
	}

	for prefix, target := range cfg.Routes {
		proxy, err := reverseProxy(prefix, target)
		if err != nil {
			return nil, err
		}

		r.Mount(prefix, mountProxy(prefix, proxy, tokenAuth, rbac[prefix]))
	}

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	return &Router{handler: r}, nil
}

func (r *Router) Handler() http.Handler { return r.handler }

func mountProxy(prefix string, proxy http.Handler, tokenAuth *jwtauth.JWTAuth, allowed map[string]struct{}) http.Handler {
	r := chi.NewRouter()
	if prefix != "/auth" {
		r.Use(middleware.JWTAuth(tokenAuth))
		r.Use(middleware.RBAC(allowed))
	}

	r.Handle("/*", proxy)
	r.Handle("/", proxy)

	return r
}

func reverseProxy(prefix, target string) (*httputil.ReverseProxy, error) {
	targetURL, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("route %s parse URL: %w", prefix, err)
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, proxyErr error) {
		log.Printf("proxy error path=%s route=%s err=%v", r.URL.Path, prefix, proxyErr)
		writeJSONError(w, http.StatusBadGateway, "bad gateway")
	}

	return proxy, nil
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]any{
			"code":    status,
			"message": message,
		},
	})
}
