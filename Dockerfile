FROM golang:alpine as builder

WORKDIR /go/src/sms_server
ENV GOPROXY=https://goproxy.cn
ENV GO111MODULE=on
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
  && apk update \
  && apk add git \
  && apk add gcc \
  && apk add libc-dev

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN  GOOS=linux GOARCH=amd64 go build -ldflags "-extldflags -static  -X 'main.Buildstamp=`date -u '+%Y-%m-%d %I:%M:%S%p'`' -X 'main.Githash=`git rev-parse HEAD`' -X 'main.Goversion=`go version`'" -o /sms_server

FROM alpine

WORKDIR /opt/

RUN apk add tzdata ca-certificates && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata && rm -rf /var/cache/apk/*

ENV check.code.time.out=60 \
    yun.pian.appkey="2103e65e0f95605fe3a896a309e1fc4e" \
    mysql.datasource.url="root:mysql@tcp(192.168.23.41:3306)/wdgl?charset=utf8mb4&parseTime=true&loc=Local"

COPY --from=builder /sms_server  .

EXPOSE 8080 8080
ENTRYPOINT ["/opt/sms_server"]