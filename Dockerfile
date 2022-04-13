FROM golang:latest

RUN mkdir /build
WORKDIR /build

RUN export GO111MODULE=on

# RUN cd /build && git clone https://github.com/bepro9/go-api.git
COPY . ./

RUN go mod tidy && go build -o /docker-gs-ping

EXPOSE 4000

ENTRYPOINT [ "/docker-gs-ping" ]
