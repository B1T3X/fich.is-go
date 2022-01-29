FROM golang:alpine3.15 as build

RUN apk update

RUN apk add -U --no-cache ca-certificates

RUN update-ca-certificates

WORKDIR /app

COPY ./src /app/src

RUN go mod init fich.is/api

RUN go get github.com/gorilla/mux

RUN go get github.com/go-redis/redis/v8

RUN echo $(ls .)

RUN go build -o fichisgo ./src

FROM alpine:latest as prod

ENV FICHIS_CERTIFICATE_FILE_PATH=/mnt/tls/certificate.crt

ENV FICHIS_CERTIFICATE_KEY_PATH=/mnt/tls/private.key

ENV FICHIS_REDIS_HOST="redis"

ENV FICHIS_REDIS_PORT=6379

ENV FICHIS_HTTPS_PORT=443

ENV FICHIS_HTTP_PORT=80

COPY --from=build /app/fichisgo /app/fichisgo

COPY --from=build etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

RUN apk add libc6-compat

CMD /app/fichisgo