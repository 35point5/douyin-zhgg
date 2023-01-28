FROM golang:latest

ENV GOPROXY https://goproxy.cn,direct
WORKDIR $GOPATH/src
COPY ./src $GOPATH/src
RUN go build -o douyin-service ./app

EXPOSE 1897
ENTRYPOINT ["./douyin-service"]