FROM golang:1.22-alpine

RUN adduser -D runner
WORKDIR /runner

COPY run.sh /run.sh
RUN chmod +x /run.sh && chown runner:runner /run.sh

# `go run` needs a writable build cache. Point it somewhere the unprivileged
# runner owns, then pre-compile the stdlib at build time so the first user
# submission is not charged for it.
ENV HOME=/home/runner \
    GOCACHE=/home/runner/.cache/go-build \
    GOPATH=/home/runner/go \
    GOFLAGS=-mod=mod \
    CGO_ENABLED=0

RUN mkdir -p "$GOCACHE" "$GOPATH" && chown -R runner:runner /home/runner

USER runner

# Warm the cache as the runner user; without this every run recompiles the
# standard library inside a network-less container and blows the time limit.
RUN printf 'package main\n\nimport (\n\t"bufio"\n\t"encoding/csv"\n\t"encoding/json"\n\t"fmt"\n\t"os"\n\t"sort"\n\t"strconv"\n\t"strings"\n\t"sync"\n\t"time"\n)\n\nfunc main() {\n\t_ = bufio.NewReader\n\t_ = csv.NewReader\n\t_ = json.Marshal\n\t_ = fmt.Println\n\t_ = os.Stdin\n\t_ = sort.Ints\n\t_ = strconv.Atoi\n\t_ = strings.TrimSpace\n\t_ = sync.WaitGroup{}\n\t_ = time.Now\n}\n' > /tmp/warm.go \
    && go build -o /dev/null /tmp/warm.go \
    && rm /tmp/warm.go

ENTRYPOINT ["/run.sh"]
