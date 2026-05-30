FROM golang:1.22-alpine AS builder
WORKDIR /workspace
COPY go.work go.work
COPY proto/ proto/
COPY pkg/ pkg/
COPY services/execution-service/ services/execution-service/
RUN cd services/execution-service && go build -ldflags="-s -w" -o /execution-service ./cmd/...

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /execution-service /execution-service
EXPOSE 50052
ENTRYPOINT ["/execution-service"]
