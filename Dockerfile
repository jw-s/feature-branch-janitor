From alpine:3.6

RUN apk update && \
apk add ca-certificates && \
rm -rf /var/cache/apk/*

COPY ./bin/janitor /janitor

ENTRYPOINT ["/janitor", "--all-namespaces"]
