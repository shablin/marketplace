# Database Migrations

> [!WARNING]
> Provided migrations implemented for MVP

> [!IMPORTANT]
> Migrations organized for using `golang-migrate` CLI tool ([learn more](https://github.com/golang-migrate/migrate))

## Format

- `<version>_<name>.up.sql`
- `<version>_<name>.down.sql`

`<version>` is 6-digit sequential number per service (`000001`, `000002`, ...)

## Directories

That's where migration SQL files stored:
- For example: `services/user-service/migrations`
- `services/<service>/migrations`

## Order (Local)

Recommended order in local env:

1. `user-service`
2. `catalog-service`
3. `cart-service`
4. `order-service`
5. `payment-service`
6. `notification-service`

Due to each service keeps separate database, this order
suits to workflow: `user > catalog > cart > order > payment > notification`

## Examples

> [!IMPORTANT]
> Replace DSN with your local database creds

### Apply

```bash
migrate -path services/user-service/migrations -database "$USER_DB_DSN" up
```

### Rollback

```bash
migrate -path services/user-service/migrations -database "$USER_DB_DSN" down 1
```