package gateway

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	middleware "github.com/shablin/marketplace/pkg/middleware"
)

type Router struct {
	handler http.Handler
}

var requestCounter uint64

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriterHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func NewRouter(cfg Config) (*Router, error) {
	logger := slog.New(
		slog.NewJSONHandler(
			cfg.LogOutput,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		),
	)

	r := chi.NewRouter()
	//r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	//r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(requestContextMiddleware(logger))

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
		proxy, err := reverseProxy(prefix, target, logger)
		if err != nil {
			return nil, err
		}

		r.Mount(prefix, mountProxy(prefix, proxy, tokenAuth, rbac[prefix]))
	}

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"status":  "ok",
			"service": "api-gateway",
		})
	})

	r.Get("/ready", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"status":  "ready",
			"service": "api-gateway",
		})
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

func reverseProxy(prefix, target string, logger *slog.Logger) (*httputil.ReverseProxy, error) {
	targetURL, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("route %s parse URL: %w", prefix, err)
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, proxyErr error) {
		logger.Error("proxy.error",
			"route", prefix,
			"path", r.URL.Path,
			"error", proxyErr.Error(),
		)
		//log.Printf("proxy error path=%s route=%s err=%v", r.URL.Path, prefix, proxyErr)
		writeJSONError(w, http.StatusBadGateway, "bad gateway")
	}

	return proxy, nil
}

func requestContextMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = generateRequestID()
			}

			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			rec.Header().Set("X-Request-ID", requestID)
			start := time.Now()
			next.ServeHTTP(rec, r)

			logger.Info("request.completed",
				"service", "api-gateway",
				"request_id", requestID,
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"duration_ms", time.Since(start).Milliseconds(),
			)
		})
	}
}

func generateRequestID() string {
	unixtime := time.Now().UnixNano()
	reqCounter := atomic.AddUint64(&requestCounter, 1)
	return "request-" +
		strconv.FormatInt(unixtime, 36) + "-" +
		strconv.FormatUint(reqCounter, 36)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{
		"error": map[string]any{
			"code":     status,
			"messsage": message,
		},
	})
}
