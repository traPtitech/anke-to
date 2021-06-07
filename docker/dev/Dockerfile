FROM golang:1.16.5-alpine as build-step
RUN apk add --update --no-cache ca-certificates git

WORKDIR /go/src/github.com/traPtitech/anke-to

RUN go get -u github.com/pilu/fresh

COPY go.mod go.sum ./
RUN go mod download

ENTRYPOINT fresh -c docker/dev/fresh.conf
