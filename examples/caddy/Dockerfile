FROM golang:1.25.0-alpine AS builder

RUN cd /tmp && \
	go run github.com/caddyserver/xcaddy/cmd/xcaddy@latest build --with github.com/dotvezz/caddy-mirror

FROM scratch

COPY --from=builder /tmp/caddy /caddy

ENTRYPOINT [ "/caddy", "run", "--config", "/etc/caddy/Caddyfile" ]
