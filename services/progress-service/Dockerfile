FROM golang:1.22-alpine AS builder
WORKDIR /workspace
COPY go.work go.work
COPY proto/ proto/
COPY pkg/ pkg/
COPY services/progress-service/ services/progress-service/
RUN cd services/progress-service && go build -ldflags="-s -w" -o /progress-service ./cmd/...

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /progress-service /progress-service
EXPOSE 50054
ENTRYPOINT ["/progress-service"]
