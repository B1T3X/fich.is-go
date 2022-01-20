FROM golang:latest as build

ENV GOPATH /go

ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

ENV FICHIS_HTTPS_PORT=443

ENV FICHIS_CERTIFICATE_FILE_PATH="/app/tls/certificate.crt"

ENV FICHIS_PRIVATE_KEY_FILE_PATH="/app/tls/private.key"

WORKDIR /app

COPY ./src /app/

RUN go get github.com/gorilla/mux

RUN go build -o fichisgo

FROM alpine:latest as prod

COPY --from=build /app/fichisgo /app

CMD ./fichisgo