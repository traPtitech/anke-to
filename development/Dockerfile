FROM golang:1.12.5-alpine as build-step
RUN apk add --update --no-cache ca-certificates git

WORKDIR /go/src/git.trapti.tech/SysAd/anke-to

RUN apk add --update --no-cache git \
  &&  go get -u github.com/golang/dep/cmd/dep \
  &&  go get -u github.com/pilu/fresh

COPY Gopkg.lock Gopkg.toml ./
RUN dep ensure --vendor-only

COPY . .

COPY development/fresh.conf ./

ENTRYPOINT fresh -c development/fresh.conf
