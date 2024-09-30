# PasteBin

一个工具，可以用来生成纯文本网页分享文字内容，或者也可以上传附件，作为一个临时的文件服务器分享给别人

单次分享的限制为 100mb，暂不支持反向代理到域名子目录，数据和日志位于 webnote_data 目录中

### 快速上手

直接运行独立的二进制文件即可，默认监听 10002 端口

```sh
./pastebin
```

### 使用说明

命令行可以接收的参数

参数|默认值|描述
-|-|-
-port|10002|程序监听的端口号

### API

- ___POST /___

请求：multipart/form-data，存储了数据内容的链接

body：`f` 文件

- ___GET /{uid}___

返回该链接所对应的内容

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
- [x] 自定义端口号
- [ ] 美化页面，重写糟糕的 js

## 许可证

MIT © ZShab Niba
