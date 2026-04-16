FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/catalog-service ./cmd/catalog-service

FROM alpine:3.23
WORKDIR /app
RUN apk add --no-cache ca-certificates wget

COPY --from=builder /bin/catalog-service /usr/local/bin/catalog-service

ENV SERVICE_NAME=catalog-service
ENV PORT=8083
EXPOSE 8083

HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=5 \
CMD wget -qO- http://127.0.0.1:${PORT}/health || exit 1

ENTRYPOINT [ "/usr/local/bin/catalog-service" ]
