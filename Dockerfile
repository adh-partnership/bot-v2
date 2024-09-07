ARG BASE_DISTRIBUTION=distroless

FROM golang:1.23 as builder

WORKDIR /work
COPY . .
RUN make build

FROM cgr.dev/chainguard/static:latest as distroless
FROM ubuntu:20.04 as debug

FROM ${BASE_DISTRIBUTION:-distroless}
COPY --from=builder /work/out/bot /usr/local/bin/bot
COPY facilities/ /etc/bot/facilities
RUN touch /etc/bot/config.yaml

ENTRYPOINT [ "/usr/local/bin/bot", "--config", "/etc/bot/config.yaml", "--facilities", "/etc/bot/facilities" ]
