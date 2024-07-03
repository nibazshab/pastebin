# pastebin

粘贴板，生成一个纯文本、图片网页，限制大小 < 5 mb

## 使用说明

默认监听 10002 端口，数据以文件的形式明文储存在 pastebin 文件同目录的 tmp 目录下（需要自己手动创建）

1. 编译源代码
2. 创建 tmp 目录
3. 运行程序 `./pastebin`

__编译步骤__

```sh
git clone https://github.com/nibazshab/pastebin.git
cd pastebin
CGO_ENABLED=0 go build -ldflags="-s -w"
```

测试平台：Linux amd64

## API

> ___POST /___

参数：`f` 文件

返回存储了数据内容的链接

> ___GET /{uid}___

返回该链接所对应的文本内容/图片

## PLAN-B

- [x] 响应 `dmesg | curl -F "f=@-" host` 形式的请求
- [ ] 解决 favicon.ico 的问题

## 开源地址

https://github.com/nibazshab/pastebin

## 许可证

MIT © ZShab Niba
