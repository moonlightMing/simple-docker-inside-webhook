FROM golang:1.9.5

WORKDIR /go/src/github.com/moonlightming/simple-docker-inside-webhook

COPY . /go/src/github.com/moonlightming/simple-docker-inside-webhook

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o main


FROM alpine:3.6

EXPOSE 9375

WORKDIR /src/app/webhook

ENV TIME_ZONE="Asia/Shanghai"

RUN apk add --no-cache tzdata ca-certificates \
     && echo ${TIME_ZONE} > /etc/timezone     \
     && ln -sf /usr/share/zoneinfo/${TIME_ZONE} /etc/localtime

# 将宿主机Docker IP地址写入本地Host
RUN /sbin/ip route|awk '/default/ { print  $3,"\tdockerhost" }' >> /etc/hosts

COPY --from=0 /go/src/github.com/moonlightming/simple-docker-inside-webhook/ /src/app/webhook

CMD ["/bin/sh", "entrypoint.sh"]