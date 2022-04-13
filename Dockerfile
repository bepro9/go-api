FROM golang:latest

RUN mkdir /build
WORKDIR /build

RUN export GO111MODULE=ON

RUN cd /build && git clone https://github.com/bepro9/go-api

RUN cd /build/go-api && go build
EXPOSE 4000

ENTRYPOINT [" /build/go-api/main"]

