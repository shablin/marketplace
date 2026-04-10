# marketplace
Yet another project in educational purposes.

# My Goals
* Implement Microservice architecture
* Implement gRPC and Protobuf

and much more...

# About
A backend for mini-marketplace without overkill stack and technologies.
This aims to implement MVP as a lightweight solution.

You can get to know my plans below:

## Root Structure
```
.
api/                          External REST contracts (OpenAPI)
proto/                        Internal gRPC contracts (Protobuf)
internal/
....platform/                 Shared tech layer
services/
....auth-service/             Authentication & Authorization
....user-service/             User profile & account data
....catalog-service/          Product catalog and search metadata
....cart-service/             Shopping cart management
....order-service/            Order creation and lifecycle
....payment-service/          Payment processing orchestra
....notification-service/     User notifications
```

## Service Structure (per service)
```
services/<service>/
....cmd/<service>/main.go     Service entrypoint
....internal/                 Internal service logic
....pkg/                      Reusable external service libs
....go.mod                    Local Go module for service
```

## Platform
`internal/platform` contains configuration, logging, middleware.
Domain/bussines logic statys within each service and is not
shared through this layer.
