# Simple-docker-inside-WebHook

自用项目，部署在容器内的Docker-WebHook，与宿主机Docker进行通信，用于更新容器内Service。CI/CD中的一环。

- 代码仓库: Github
- 容器构建: 阿里云容器服务
- 部署集成: 本项目处于这一环

## 功能:
- WebHook钩子触发自动更新Service
- 更新完毕邮件通知（可选）
- 定时清理无用镜像存储

## 构建

我直接使用了阿里云的容器仓库进行自动化构建，如需本地构建，项目中也提供了Dockerfile。
项目中使用了多阶段构建，初次构建会拉取大镜像耗时较长，但构建后镜像体积很小。

```
git clone https://github.com/moonlightMing/simple-docker-inside-webhook.git
cd ./simple-docker-inside-webhook

vim config.json

# 添加如下内容, 使用时请去除注释，json语法不允许有注释
{
  # 本地监听端口
  "bind_host": ":9375",
  "docker_host": "dockerhost:2375",
  "docker_api_version": "x.xx",
  "auth_key": "xxxxxx",
  # 远程仓库账号
  "docker_registry_auth": {
    "user": "xxxxxx@xx.com",
    "password": "xxxxxxxx"
  },
  # 邮件设置
  "email": {
    # 开关
    "open": false,
    # 发送者邮箱配置
    "smtp_host": "smtp.163.com",
    "smtp_port": ":465",
    "user_email": "xxx@163.com",
    "password": "xxx",
    # 接受者地址
    "send_to": "xxx@qq.com"
  },
  # 定时任务 不用可以为空
  "cron": [
    {
      "event": "CLEAN_NONE_TAG_IMAGE",
      "spec": "* * * */7 * *"
    }
  ]
}

docker build -t xxx/docker-webhook .
```
