# pastebin

粘贴板，生成一个纯文本、图片网页，限制大小 < 5 mb

## 使用说明

默认监听 10002 端口，数据储存在 pastebin 可执行文件同级 data 目录下的 pastebin.db 文件中

1. 编译源代码
3. 运行程序 `./pastebin`

__编译步骤__

编译依赖：gcc，go

```sh
git clone https://github.com/nibazshab/pastebin.git
cd pastebin
go get ./...
CGO_ENABLED=1 go build -ldflags="-s -w"
```

> [!TIP]
> - 测试平台：Linux amd64
> - 反向代理时不要代理到域名子目录
> - 使用了 mattn/go-sqlite3 库而受限，编译产物非静态文件

__API__

> ___POST /___

请求：multipart/form-data，存储了数据内容的链接

body：`f` 文件

> ___GET /{uid}___

返回该链接所对应的文本内容/图片

## PLAN-B

- [x] 响应 `dmesg | curl -F "f=@-" host` 形式的请求
- [x] 解决 favicon.ico 的问题
- [x] 变更相对路径为绝对路径

## 许可证

MIT © ZShab Niba
