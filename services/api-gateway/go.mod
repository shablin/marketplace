module github.com/shablin/marketplace/services/api-gateway

go 1.26.1

require (
	github.com/go-chi/chi/v5 v5.2.5
	github.com/go-chi/jwtauth/v5 v5.4.0
	github.com/shablin/marketplace v0.0.0-20260410122912-ef82a847fd0b
)

require (
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.4.0 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/lestrrat-go/blackmagic v1.0.3 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/httprc/v3 v3.0.0-beta2 // indirect
	github.com/lestrrat-go/jwx/v3 v3.0.2 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
)

replace github.com/shablin/marketplace => ../..
