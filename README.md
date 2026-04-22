# marketplace

Study backend project with microservices for minimal marketplace

# About

A brief description of repository

## Structure

```
api/                        External REST contracts (OpenAPI)
proto/                      Internal gRPC contracts (Protobuf)
internal/platform/          Shared infrastructure docs

services/
...api-gateway/
...auth-service/
...cart-service/
...catalog-service/
...notification-service/
...order-service/
...payment-service/
...user-service/
...
```

## Quick Start
### Step 1. Up the services

```bash
docker compose up --build
```

By default, services available at:

- gateway:          `:8080`
- auth:             `:8081`
- cart:             `:8082`
- catalog:          `:8083`
- notification:     `:8084`
- order:            `:8085`
- payment:          `:8086`
- user:             `:8087`

### Step 2. Check services availability

```bash
curl -i http://localhost:8081/health
curl -i http://localhost:8081/ready
```

The same way for other services, just
get `/health` and `/ready`
with CURL by service's port, or run loop:

```bash
for port in 8081 8082 8083 8084 8084 8085 8086 8087; do
     curl -fsS "http://localhost:${port}/health" && echo " <- ${port} OK"
done
```

## Workflow
### Step 1. Register

```bash
curl -s -X POST http://localhost:8081/auth/register \
     -H 'Content-Type: application/json' \
     -d '{"email":"you@example.com", "password": "yourpassword"}'
```

### Step 2. Login

```bash
curl -s -X POST http://localhost:8081/auth/login \
     -H 'Content-Type: application/json' \
     -d '{"email":"you@example.com", "password":"yourpassword"}'
```

> [!WARNING]
> Save `access_token` and `refresh_token` from response

### Step 3. Checkout

Checkout pipeline consist of:

1. Create order
2. Payment
3. Get notification

You can test this pipeline with integration test:

```bash
cd services/order-service && go test ./internal/order -run TestCreateOrderPaymentNotificationFlow
```

> [!IMPORTANT]
> Now this integration test and related business logic are still not implemented.
> Reference only

### Step 4. Logout / Revoke Token

```bash
curl -s -X POST http://localhost:8081/auth/logout \
     -H 'Content-Type: application/json' \
     -d '{"refresh_token":"<refresh_token>"}'
```

### Step 5. Reset password

```bash
curl -s -X POST http://localhost:8081/auth/reset-password \
     -H 'Content-Type: application/json' \
     -d '{"email":"you@example.com"}'
```

Then you should approve password reset:

```bash
curl -s -X POST http://localhost:8081/auth/reset-password/confirm \
     -H 'Content-Type: application/json' \
     -d '{"token":"<reset_token>", "new_password":"yournewpassword"}'
```

## See also

- REST/OpenAPI: `api/openapi.yaml`
- gRPC/Protobuf: `proto/` and `proto/README.md`