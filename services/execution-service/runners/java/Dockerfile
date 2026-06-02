FROM eclipse-temurin:21-jre-alpine

RUN adduser -D runner
WORKDIR /runner
# JDK tools needed for compilation
RUN apk add --no-cache openjdk21-jdk

COPY run.sh /run.sh
RUN chmod +x /run.sh && chown runner:runner /run.sh

USER runner

ENTRYPOINT ["/run.sh"]

