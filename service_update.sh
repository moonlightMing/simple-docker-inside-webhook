#!/bin/bash

docker pull registry.cn-shenzhen.aliyuncs.com/moonlightming/simple-docker-inside-webhook:latest
docker service update --image registry.cn-shenzhen.aliyuncs.com/moonlightming/simple-docker-inside-webhook:latest webhook_simple-docker-inside-webhook --with-registry-auth