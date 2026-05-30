FROM golang:1.22-alpine AS builder
WORKDIR /workspace
COPY go.work go.work
COPY proto/ proto/
COPY pkg/ pkg/
COPY services/api-gateway/ services/api-gateway/
RUN cd services/api-gateway && go build -ldflags="-s -w" -o /api-gateway ./cmd/...

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /api-gateway /api-gateway
EXPOSE 8080
ENTRYPOINT ["/api-gateway"]
