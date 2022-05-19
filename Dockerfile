FROM golang:alpine3.15 as build

RUN apk update

RUN apk add -U --no-cache ca-certificates

RUN update-ca-certificates

WORKDIR /app

COPY ./src /app/src

RUN go mod init github.com/b1t3x/fich.is-go

RUN go get github.com/gorilla/mux

RUN go get cloud.google.com/go/firestore

RUN go get google.golang.org/api/option

RUN echo $(ls .)

RUN go build -o fichisgo ./src

FROM alpine:latest as prod

ENV FICHIS_CERTIFICATE_FILE_PATH=/mnt/tls/certificate.crt

ENV FICHIS_CERTIFICATE_KEY_PATH=/mnt/tls/private.key

ENV FICHIS_HTTPS_PORT=443

ENV FICHIS_HTTP_PORT=80

ENV FICHIS_DOMAIN_NAME="fich.is"

env FICHIS_GOOGLE_APPLICATION_CREDENTIALS_FILE_PATH="/app/secrets/sa.json"

env FICHIS_GOOGLE_PROJECT_ID="fichis-go"

ENV FICHIS_PROBE_PATH="/api/healthprobe"

COPY --from=build /app/fichisgo /app/fichisgo

COPY --from=build etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

RUN apk add libc6-compat

CMD /app/fichisgo
