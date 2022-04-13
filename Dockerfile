FROM golang:latest

RUN mkdir /build

RUN export GO111MODULE=ON

RUN cd /build && git clone https://github.com/bepro9/go-api

RUN cd go-api && go build
RUN go run main.go

EXPOSE 4000


