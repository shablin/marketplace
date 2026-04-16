FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/order-service ./cmd/order-service

FROM alpine:3.23
WORKDIR /app
RUN apk add --no-cache ca-certificates wget

COPY --from=builder /bin/order-service /usr/local/bin/order-service

ENV SERVICE_NAME=order-service
ENV PORT=8085
EXPOSE 8085

HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=5 \
CMD wget -qO- http://127.0.0.1:${PORT}/health || exit 1

ENTRYPOINT [ "/usr/local/bin/order-service" ]
