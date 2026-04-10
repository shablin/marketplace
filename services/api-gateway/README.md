# api-gateway

Public REST entrypoint

## Routes

* `/auth`
* `/users`
* `/products`
* `/cart`
* `/orders`
* `/payments`
* `/notifications`

## Features

* JWT auth
* RBAC middleware for `buyer`, `seller`, `admin`
* Reverse proxy routing
* Basic JSON error handling for proxy and panic cases
