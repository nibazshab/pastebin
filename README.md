# pastebin

粘贴板，生成一个纯文本、图片网页，限制大小 < 5 mb

## 使用说明

纯静态无任何依赖，直接运行 `pastebin` 二进制文件即可，默认监听 10002 端口，内容储存在 sqlite 数据库中

数据库文件位于文件同级目录的 data/pastebin.db，日志文件位于 data/log.log

命令行参数|默认值|描述
-|-|-
-port|10002|文件服务器监听的端口

> [!TIP]
> - 测试平台：Linux amd64
> - 暂不支持反向代理到域名子目录

### API

- ___POST /___

请求：multipart/form-data，存储了数据内容的链接

body：`f` 文件

- ___GET /{uid}___

返回该链接所对应的文本内容/图片

## 构建说明

所需软件包：go, musl

go 使用包管理器或任意方式安装，musl 可以通过如下命令安装

```sh
musl="https://musl.cc/x86_64-linux-musl-cross.tgz"
wget -O- "$musl" | tar -zxvf - --strip-components=1 -C /usr/local
```

开始构建

```sh
go get ./...

flags="-s -w --extldflags '-static -fpic'"
export GOOS=linux
export GOARCH=amd64
export CC=x86_64-linux-musl-gcc
export CGO_ENABLED=1

go build -ldflags="$flags"
```

## PLAN-B

- [x] 响应 `dmesg | curl -F "f=@-" host` 形式的请求
- [x] 解决 favicon.ico 的问题
- [x] 变更相对路径为绝对路径

## 许可证

MIT © ZShab Niba
