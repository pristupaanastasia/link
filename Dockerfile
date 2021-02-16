FROM golang:latest
RUN mkdir -p /go/src/app
WORKDIR /go/src/app
RUN apt-get -y update && apt-get -y install git
RUN go mod init github.com/my/repo
RUN go get github.com/go-redis/redis/v8
ENTRYPOINT go run server.go