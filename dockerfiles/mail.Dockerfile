FROM golang:1.24 AS build-stage

WORKDIR /workdir
COPY app/ /workdir/app
WORKDIR  /workdir/app/cmd/mail
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app


####
FROM python:latest AS template-stage
RUN pip install jinja2-cli
COPY postfix/main.cf.tmpl /main.cf.tmpl
RUN jinja2 /main.cf.tmpl           \
    -D hostname=tlsrpt.alexsci.com \
    -o /main.cf


####
FROM ubuntu:24.04

RUN --mount=target=/var/lib/apt/lists,type=cache,sharing=locked \
    --mount=target=/var/cache/apt,type=cache,sharing=locked \
    rm -f /etc/apt/apt.conf.d/docker-clean \
    && apt-get update \
    && apt-get -y install postfix

RUN useradd -m -s /bin/bash catchall

COPY postfix/entrypoint.sh /entrypoint.sh
COPY --from=template-stage main.cf /etc/postfix/main.cf
COPY --from=build-stage /app /app

EXPOSE 25/tcp 

# Save the Postgres password to somewhere the mailbox_command can access it

CMD [ "/entrypoint.sh" ]

