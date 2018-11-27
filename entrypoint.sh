#!/bin/sh

# 将宿主机Docker IP地址写入本地Host
/sbin/ip route|awk '/default/ { print  $3,"\tdockerhost" }' >> /etc/hosts

./main