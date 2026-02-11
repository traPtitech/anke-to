FROM golang:1.25.7-alpine as builder
RUN apk add --update --no-cache ca-certificates git \
 && apk add --no-cache gcc libc-dev \
 && apk add --no-cache openssl

ENV DOCKERIZE_VERSION v0.10.0
RUN wget https://github.com/jwilder/dockerize/releases/download/${DOCKERIZE_VERSION}/dockerize-alpine-linux-amd64-${DOCKERIZE_VERSION}.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-${DOCKERIZE_VERSION}.tar.gz \
    && rm dockerize-alpine-linux-amd64-${DOCKERIZE_VERSION}.tar.gz

WORKDIR /go/src/github.com/traPtitech/anke-to


COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install go.uber.org/mock/mockgen@v0.6.0  \
 && go generate ./model ./traq

ENTRYPOINT dockerize -timeout 10s -wait tcp://mysql:3306 sh -c "go test ./model -v -cover && go test ./controller -v -cover"
