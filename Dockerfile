FROM golang:1.12beta1-alpine3.8

ENV GO111MODULE on

ARG pkg=github.com/kutsuzawa/slackbot

RUN apk add --no-cache ca-certificates

COPY . $GOPATH/src/$pkg

RUN set -ex \
      && apk add --no-cache --virtual git \
      && cd $GOPATH/src/$pkg \
      && go build \
      && go install

COPY .env /go
CMD ["slackbot"]
