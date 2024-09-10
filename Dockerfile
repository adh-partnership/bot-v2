ARG BASE_DISTRIBUTION=distroless

FROM golang:1.23 as builder

WORKDIR /work
COPY . .
RUN make build
RUN touch /work/out/config.yaml

FROM cgr.dev/chainguard/static:latest as distroless
FROM ubuntu:20.04 as debug

FROM ${BASE_DISTRIBUTION:-distroless}
COPY --from=builder /work/out/bot /usr/local/bin/bot
COPY --from=builder /work/out/config.yaml /etc/bot/config.yaml
COPY facilities/ /etc/bot/facilities

ENTRYPOINT [ "/usr/local/bin/bot", "start", "--config", "/etc/bot/config.yaml", "--facility-configs-path", "/etc/bot/facilities" ]
