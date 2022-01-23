FROM golang:alpine3.15 as build

RUN apk update

RUN apk add -U --no-cache ca-certificates

RUN update-ca-certificates

WORKDIR /app

COPY ./src /app

RUN go get github.com/gorilla/mux

RUN go build -o fichisgo .

FROM alpine:latest as prod

ENV FICHIS_CERTIFICATE_FILE_PATH=/mnt/tls/certificate.crt

ENV FICHIS_CERTIFICATE_KEY_PATH=/mnt/tls/private.key

ENV FICHIS_REDIS_HOST="redis"

ENV FICHIS_REDIS_PORT=6379

ENV FICHIS_HTTPS_PORT=443

COPY --from=build /app/fichisgo /app/fichisgo

COPY --from=build etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

RUN apk add libc6-compat

CMD /app/fichisgo