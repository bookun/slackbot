FROM golang

ENV GO111MODULE on

RUN mkdir -p ${GOPATH}/src/github.com/kutsuzawa
WORKDIR ${GOPATH}/src/github.com/kutsuzawa
RUN git clone https://github.com/kutsuzawa/slackbot.git
WORKDIR ${GOPATH}/src/github.com/kutsuzawa/slackbot
RUN cp .env.sample .env

RUN go build ./...
RUN go install

ENTRYPOINT ["slackbot"]

EXPOSE 3030
