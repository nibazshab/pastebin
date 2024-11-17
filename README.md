# PasteBin

一个工具，可以用来生成纯文本网页分享文字内容，或者也可以上传附件，作为一个临时的文件服务器分享给别人

单次分享的限制为 100mb，暂不支持反向代理到域名子目录，数据和日志默认位于 pastebin_data 目录中

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
-path|pastebin_data|数据目录（相对程序文件的路径）

### API

- ___POST /___

请求：multipart/form-data，返回存储了数据内容的链接

body：`f=@文件` 上传的文件，（可选）`t=text/file` 文件类型，（可选）`v=true/false` 是否可预览

- ___GET /{uid}___

返回该链接所对应的内容

## 构建说明

所需软件包：go, musl

go 使用包管理器或任意方式安装，musl 可以通过如下命令安装

```sh
wget -O musl.tgz https://musl.cc/x86_64-linux-musl-cross.tgz
tar -zxvf musl.tgz --strip-components=1 -C /usr/local
```

开始构建

```sh
go mod tidy
flags="-s -w --extldflags '-static' \
 -X main.version=$(git describe --abbrev=0 --tags)"
export CC=x86_64-linux-musl-gcc
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=1
go build -ldflags="$flags"
```

## 打包软件包

PKGBUILD

```sh
pkgname=pastebin
pkgver=
pkgrel=1
arch=('x86_64')
source=("$pkgname" "$pkgname.service")
sha256sums=('SKIP' 'SKIP')
package() {
  mkdir -p $pkgdir/usr/local/$pkgname
  cp $pkgname $pkgdir/usr/local/$pkgname
  mkdir -p $pkgdir/usr/lib/systemd/system
  cp $pkgname.service $pkgdir/usr/lib/systemd/system
}
```

pastebin.service

```ini
[Unit]
Description=pastebin service
[Service]
ExecStart=/usr/local/pastebin/pastebin
Restart=on-failure
[Install]
WantedBy=multi-user.target
```

## PLAN-B

- [x] 响应 `dmesg | curl -F "f=@-" host` 形式的请求
- [x] 解决 favicon.ico 的问题
- [x] 变更相对路径为绝对路径
- [x] 自定义端口号
- [x] 美化页面，重写糟糕的 js
- [x] 修复 xlsx, word 等被检测为 zip 类型的问题（直接返回原始文件名）
- [x] ctrl v 上传图片
- [ ] 支持绝对路径的数据目录

## 许可证

MIT © ZShab Niba
