FROM golang:alpine3.15 as build

RUN apk update

RUN apk add -U --no-cache ca-certificates && update-ca-certificates

WORKDIR /app

COPY ./src /app

RUN go get github.com/gorilla/mux

RUN go build -o fichisgo

FROM alpine:latest as prod

ENV GOPATH /go

ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

ENV FICHIS_HTTPS_PORT=443

ENV FICHIS_CERTIFICATE_FILE_PATH="/app/tls/certificate.crt"

ENV FICHIS_PRIVATE_KEY_FILE_PATH="/app/tls/private.key"

COPY --from=build /app/fichisgo /app/fichisgo

COPY --from=build etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

RUN update-ca-certificates

RUN apk add libc6-compat

CMD /app/fichisgo