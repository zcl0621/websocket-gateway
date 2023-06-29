FROM golang:1.18 AS build-dist
ENV GOSUMDB='sum.golang.google.cn'
ENV GOPROXY='http://nexus.prd.estargo.com.cn:8081/repository/ustc-go-proxy-repo/,direct'

RUN sed -i 's@deb.debian.org@nexus.prd.estargo.com.cn:8081/repository/ustc-debian-proxy-repo@g' /etc/apt/sources.list

RUN apt-get install -y git

WORKDIR /go/cache

ADD go.mod .
ADD go.sum .
RUN go mod download

WORKDIR /go/release

ADD . .
RUN go mod tidy
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -tags netgo -installsuffix cgo -o /bin/app main.go

FROM alpine as prod

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

RUN apk add --no-cache -U  tzdata

RUN  ln -sf /usr/share/zoneinfo/Asia/Shanghai  /etc/localtime && \
     echo "Asia/Shanghai" > /etc/timezone

COPY --from=build-dist /bin/app /bin/app
RUN chmod +x /bin/app

ENV RUN_MODE='release'
CMD ["/bin/app"]
EXPOSE 8080
