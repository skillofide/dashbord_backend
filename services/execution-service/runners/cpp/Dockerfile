FROM alpine:3.20

RUN apk add --no-cache g++ && adduser -D runner
WORKDIR /runner
COPY run.sh /run.sh
RUN chmod +x /run.sh && chown runner:runner /run.sh

USER runner

ENTRYPOINT ["/run.sh"]
