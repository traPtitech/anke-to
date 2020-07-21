# build backend
FROM golang:1.12.5-alpine as server-build
RUN apk add --update --no-cache ca-certificates git

WORKDIR /github.com/traPtitech/anke-to

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /anke-to

#build frontend
FROM node:12-alpine as client-build
WORKDIR /github.com/traPtitech/anke-to/client
COPY ./client/package*.json ./
RUN npm ci
COPY ./client .
RUN npm run build


# run

FROM alpine:3.12.0
WORKDIR /app

ENV DOCKERIZE_VERSION v0.6.1

RUN apk --update add tzdata \
  && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
  && wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
  && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
  && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
  && apk add --update ca-certificates \
  && update-ca-certificates \
  && rm -rf /var/cache/apk/*

COPY --from=server-build /anke-to ./
COPY --from=client-build /github.com/traPtitech/anke-to/client/dist ./client/dist/
ENTRYPOINT ./anke-to
