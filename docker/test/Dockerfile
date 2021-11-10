FROM golang:1.17.3-alpine as builder
RUN apk add --update --no-cache ca-certificates git
RUN apk add --no-cache gcc libc-dev
RUN apk add --no-cache openssl

ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/${DOCKERIZE_VERSION}/dockerize-alpine-linux-amd64-${DOCKERIZE_VERSION}.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-${DOCKERIZE_VERSION}.tar.gz \
    && rm dockerize-alpine-linux-amd64-${DOCKERIZE_VERSION}.tar.gz

WORKDIR /go/src/github.com/traPtitech/anke-to


COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/golang/mock/mockgen  \
 && go generate ./model ./traq

ENTRYPOINT dockerize -timeout 10s -wait tcp://mysql:3306 go test ./model ./router -v cover

