FROM postgres:16-alpine

# python3 runs the test harness; no pip packages are needed because queries are
# executed through psql rather than a database driver.
RUN apk add --no-cache python3

# Initialise the cluster at BUILD time so container start only has to boot an
# already-created database. initdb takes over a second, which would otherwise be
# paid on every single test case.
ENV PGDATA=/var/lib/postgresql/judge
ENV POSTGRES_HOST_AUTH_METHOD=trust

USER postgres
RUN mkdir -p "$PGDATA" \
 && initdb -D "$PGDATA" --auth=trust --username=postgres --encoding=UTF8 \
 # fsync off and a small shared buffer: this cluster is thrown away after one
 # test case, so durability is worthless and startup speed is everything.
 && echo "fsync = off"                      >> "$PGDATA/postgresql.conf" \
 && echo "synchronous_commit = off"          >> "$PGDATA/postgresql.conf" \
 && echo "full_page_writes = off"            >> "$PGDATA/postgresql.conf" \
 && echo "listen_addresses = ''"             >> "$PGDATA/postgresql.conf" \
 && echo "unix_socket_directories = '/tmp'"  >> "$PGDATA/postgresql.conf" \
 && pg_ctl -D "$PGDATA" -o "-c listen_addresses=''" -w start \
 && createdb -h /tmp judge \
 && pg_ctl -D "$PGDATA" -m fast -w stop

ENV PGHOST=/tmp
ENV PGUSER=postgres

COPY harness.py /harness.py
COPY run.sh /run.sh

USER root
RUN chmod +x /run.sh /harness.py
USER postgres

ENTRYPOINT ["/run.sh"]
