# Simple-docker-inside-WebHook

自用项目，部署在容器内的WebHook，与宿主机Docker进行通信，用于更新容器内Service。CI/CD中的一环。

- 代码仓库: Github
- 容器构建: 阿里云容器服务
- 部署集成: 本项目处于这一环

## 功能:
- WebHook自动更新Service
- 更新完毕邮件通知
- 定时清理无用镜像存储