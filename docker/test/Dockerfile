FROM golang:1.17.3-alpine as builder
RUN apk add --update --no-cache ca-certificates git
RUN apk add --no-cache gcc libc-dev

WORKDIR /go/src/github.com/traPtitech/anke-to


COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/golang/mock/mockgen  \
 && go generate ./model ./traq

CMD go test ./model ./router -v cover

