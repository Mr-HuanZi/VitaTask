# VitaTask

这是一款轻量级的在线项目任务管理工具，提供项目管理、任务分发、即时IM等工具。

> [UI仓库](https://github.com/Mr-HuanZi/VitaTask-UI)

## 环境准备

- 要求`Go1.20+`，需要开启Go Mod
- Mysql5.7

### 依赖

```shell
go mod tidy
```

### Mysql

安装MySQL数据(建议为5.7版本)并创建好名为`vita_task`的数据库，默认字符集为`utf8-mb4`。

> 如需要修改请在`app.yml`内指定

将`vita_task.sql`导入到Mysql中。

## 运行

API服务

```shell
go run ./cmd/vita-api/main.go
```

Websocket Gateway服务

```shell
go run ./cmd/vita-gateway/main.go
```

## 配置文件

执行以下命令将配置文件拷贝一份:

```shell
cp ./config/app.example.yaml ./config/app.yml
```

### 配置项说明

#### app

设置程序运行必要配置，如端口号、调试模式。

#### mysql

连接MySQL数据库配置

#### redis

连接Redis配置，暂时无用

#### jwt

JWT配置项目，用于设置加密Key、过期时间、签名信息。

#### member

设置成员默认密码。

#### gateway

websocket网关配置