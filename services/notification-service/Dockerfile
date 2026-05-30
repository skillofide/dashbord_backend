FROM golang:1.22-alpine AS builder
WORKDIR /workspace
COPY go.work go.work
COPY proto/ proto/
COPY pkg/ pkg/
COPY services/notification-service/ services/notification-service/
RUN cd services/notification-service && go build -ldflags="-s -w" -o /notification-service ./cmd/...

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /notification-service /notification-service
EXPOSE 8081
ENTRYPOINT ["/notification-service"]
