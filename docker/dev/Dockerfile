FROM golang:1.26.5-alpine@sha256:0178a641fbb4858c5f1b48e34bdaabe0350a330a1b1149aabd498d0699ff5fb2 as build-step
RUN apk add --update --no-cache ca-certificates git

WORKDIR /go/src/github.com/traPtitech/anke-to

RUN go install github.com/cosmtrek/air@v1.64.5

COPY go.mod go.sum ./
RUN go mod download

ENV DOCKERIZE_VERSION v0.6.1
RUN apk add --no-cache openssl \
 && wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
 && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
 && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

ENTRYPOINT dockerize -timeout 10s -wait tcp://mysql:3306  air -c docker/dev/.air.toml