FROM python:3.12-slim

# Install no extra packages to minimize attack surface
RUN useradd --create-home --shell /bin/bash runner

WORKDIR /runner
COPY run.sh /run.sh
RUN chmod +x /run.sh && chown runner:runner /run.sh

USER runner

ENTRYPOINT ["/run.sh"]
