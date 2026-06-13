FROM golang:1.26.4-alpine@sha256:7a3e50096189ad57c9f9f865e7e4aa8585ed1585248513dc5cda498e2f41812c as builder
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
