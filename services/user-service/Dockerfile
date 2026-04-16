FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/user-service ./cmd/user-service

FROM alpine:3.23
WORKDIR /app
RUN apk add --no-cache ca-certificates wget

COPY --from=builder /bin/user-service /usr/local/bin/user-service

ENV SERVICE_NAME=user-service
ENV PORT=8082
EXPOSE 8082

HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=5 \
CMD wget -qO- http://127.0.0.1:${PORT}/health || exit 1

ENTRYPOINT [ "/usr/local/bin/user-service" ]
