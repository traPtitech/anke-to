# build
FROM golang:1.12.5-alpine as build
RUN apk add --update --no-cache ca-certificates git

WORKDIR /git.trap.jp/SysAd/
# githubへの移行後
# WORKDIR /go/src/github.com/traPtitech/anke-to

RUN apk add --update --no-cache git \
  &&  go get -u github.com/golang/dep/cmd/dep

COPY Gopkg.lock Gopkg.toml ./
RUN dep ensure --vendor-only

COPY . .

RUN go build -o /anke-to

# run

FROM alpine:3.9
WORKDIR /app

ENV DOCKERIZE_VERSION v0.6.1

RUN apk --update add tzdata \
  && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
  && wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
  && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
  && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

COPY --from=build /anke-to ./