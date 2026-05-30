FROM golang:1.22-alpine AS builder
WORKDIR /workspace
COPY go.work go.work
COPY proto/ proto/
COPY pkg/ pkg/
COPY services/submission-service/ services/submission-service/
RUN cd services/submission-service && go build -ldflags="-s -w" -o /submission-service ./cmd/...

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /submission-service /submission-service
EXPOSE 50053
ENTRYPOINT ["/submission-service"]
