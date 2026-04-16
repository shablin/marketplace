package main

import (
	"log"
	"net/http"

	"github.com/shablin/marketplace/services/api-gateway/internal/gateway"
)

func main() {
	cfg := gateway.NewConfig()
	router, err := gateway.NewRouter(cfg)
	if err != nil {
		log.Fatalf("build gateway router: %v", err)
	}

	log.Printf("api-gateway listening on %s", cfg.ListenAddr)
	if err := http.ListenAndServe(cfg.ListenAddr, router.Handler()); err != nil {
		log.Fatalf("api-gateway stopped: %v", err)
	}
}
