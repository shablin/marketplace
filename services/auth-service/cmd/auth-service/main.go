package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}

	return fallback
}

func main() {
	serviceName := getenv("SERVICE_NAME", "auth-service")
	port := getenv("PORT", "8081")

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "%s is running", serviceName)
	})

	addr := fmt.Sprintf(":%s", port)
	log.Printf("%s starting on %s", serviceName, addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("%s stopped: %v", serviceName, err)
	}
}
