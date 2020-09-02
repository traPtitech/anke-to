FROM golang:1.15.1-alpine as build-step
RUN apk add --update --no-cache ca-certificates git

WORKDIR /go/src/github.com/traPtitech/anke-to

RUN apk add --update --no-cache git \
  &&  go get -u github.com/pilu/fresh

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY development/fresh.conf ./

ENTRYPOINT fresh -c development/fresh.conf
