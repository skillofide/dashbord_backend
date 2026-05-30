FROM eclipse-temurin:21-jre-alpine

RUN adduser -D runner
WORKDIR /runner
COPY run.sh /run.sh
RUN chmod +x /run.sh && chown runner:runner /run.sh

# JDK tools needed for compilation
RUN apk add --no-cache openjdk21-jdk

USER runner

ENTRYPOINT ["/run.sh"]
