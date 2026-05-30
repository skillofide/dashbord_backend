FROM gcc:13-alpine

RUN adduser -D runner
WORKDIR /runner
COPY run.sh /run.sh
RUN chmod +x /run.sh && chown runner:runner /run.sh

USER runner

ENTRYPOINT ["/run.sh"]
