# syntax = docker/dockerfile:1.0-experimental

# build backend
FROM golang:1.15.3-alpine as server-build
RUN --mount=type=cache,id=aptgit,target=/var/cache/apt --mount=type=cache,target=/var/lib/apt \
  apk add --update git

WORKDIR /github.com/traPtitech/anke-to

COPY go.mod go.sum ./
RUN --mount=type=cache,id=gomod,target=/go/pkg/mod/cache \
  go mod download

COPY . .

RUN --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
  go build -o /anke-to -ldflags "-s -w"

#build frontend
FROM node:12-alpine as client-build
WORKDIR /github.com/traPtitech/anke-to/client
COPY client/package.json client/package-lock.json ./
RUN --mount=type=cache,id=npmci,target=/root/.npm \
  npm ci --prefer-offline
COPY client .
RUN --mount=type=cache,id=vuebuild,target=/github.com/traPtitech/anke-to/client/node_modules/.cache \
  npm run staging-build


# run
FROM alpine:3.12.0
WORKDIR /app

RUN apk --update --no-cache add tzdata \
  && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
  && apk del tzdata \
  && mkdir -p /usr/share/zoneinfo/Asia \
  && ln -s /etc/localtime /usr/share/zoneinfo/Asia/Tokyo
RUN apk --update --no-cache add ca-certificates \
  && update-ca-certificates \
  && rm -rf /usr/share/ca-certificates /etc/ssl/certs

COPY --from=server-build /anke-to ./
COPY --from=client-build /github.com/traPtitech/anke-to/client/dist ./client/dist/
ENTRYPOINT ./anke-to
