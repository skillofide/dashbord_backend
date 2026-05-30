FROM golang:1.22-alpine AS builder
WORKDIR /workspace

# Copy workspace files
COPY go.work go.work
COPY proto/ proto/
COPY pkg/ pkg/
COPY services/problem-service/ services/problem-service/

RUN cd services/problem-service && go build -ldflags="-s -w" -o /problem-service ./cmd/...

# ── Runtime image ─────────────────────────────────────────────────────────────
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /problem-service /problem-service
EXPOSE 50051
ENTRYPOINT ["/problem-service"]
