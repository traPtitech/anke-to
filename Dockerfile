# build backend
FROM golang:1.15.2-alpine as server-build
RUN apk add --update --no-cache git

WORKDIR /github.com/traPtitech/anke-to

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /anke-to -ldflags "-s -w"

#build frontend
FROM node:12-alpine as client-build
WORKDIR /github.com/traPtitech/anke-to/client
COPY client/package.json client/package-lock.json ./
RUN npm ci
COPY client .
RUN npm run build


# run
FROM alpine:3.12.0
WORKDIR /app

RUN apk --update --no-cache add tzdata \
  && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
  && apk del tzdata
RUN apk add --no-cache --update ca-certificates \
  && update-ca-certificates

COPY --from=server-build /anke-to ./
COPY --from=client-build /github.com/traPtitech/anke-to/client/dist ./client/dist/
ENTRYPOINT ./anke-to
