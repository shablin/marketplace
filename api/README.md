# API Contracts

This dir stores external REST contracts (OpenAPI specs)

Base HTTP API contract in `openapi.yaml` for
health checks and auth flows:

|       Flow        |           Endpoints       |
|       ---         |               ---         |
| Health Checks     | `/health` <br> `/ready`   |
| Auth              | `/auth/register` <br> `/auth/login` <br> `/auth/refresh` <br> `/auth/logout` <br> `/auth/reset-password` <br> `/auth/reset-password/confirm` |