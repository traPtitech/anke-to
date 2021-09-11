# build backend
FROM golang:1.17.1-alpine as server-build
RUN apk add --update --no-cache git

WORKDIR /github.com/traPtitech/anke-to

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /anke-to -ldflags "-s -w"

#build frontend
FROM node:14-alpine as client-build
WORKDIR /github.com/traPtitech/anke-to/client
RUN apk add --update --no-cache python3 make g++
COPY client/package.json client/package-lock.json ./
RUN npm ci
RUN npx browserslist@latest --update-db
COPY client .
RUN npm run build


# run
FROM alpine:3.14.2
WORKDIR /app

RUN apk --update --no-cache add tzdata \
  && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
  && apk del tzdata \
  && mkdir -p /usr/share/zoneinfo/Asia \
  && ln -s /etc/localtime /usr/share/zoneinfo/Asia/Tokyo
RUN apk --update --no-cache add ca-certificates \
  && update-ca-certificates \
  && rm -rf /usr/share/ca-certificates

COPY --from=server-build /anke-to ./
COPY --from=client-build /github.com/traPtitech/anke-to/client/dist ./client/dist/
ENTRYPOINT ./anke-to
