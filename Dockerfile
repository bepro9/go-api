FROM golang:latest

RUN mkdir /build
WORKDIR /build

RUN export GO111MODULE=on

COPY . ./

RUN go build -o /docker-gs-ping

EXPOSE 4000

ENTRYPOINT [ "/docker-gs-ping" ]
