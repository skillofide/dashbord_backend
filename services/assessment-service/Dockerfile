FROM golang:1.22-alpine AS builder
WORKDIR /workspace

# Copy workspace files
COPY go.work go.work
COPY go.work.sum go.work.sum
COPY proto/ proto/
COPY pkg/ pkg/
COPY services/ services/
COPY vendor/ vendor/
RUN go build -mod=vendor -ldflags="-s -w" -o /assessment-service ./services/assessment-service/cmd/...

# ── Runtime image ─────────────────────────────────────────────────────────────
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /assessment-service /assessment-service
EXPOSE 50056
ENTRYPOINT ["/assessment-service"]
