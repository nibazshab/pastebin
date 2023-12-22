# pastebin

文本粘贴板，生成一个纯文本的网页

附带一个浏览器前端页面

## 使用说明

默认监听 10002 端口，数据以文件的形式明文储存在 pastebin 文件同目录的 tmp 目录下（需要自己手动创建）

1. 编译 main.go
2. 创建 tmp 目录
3. 运行程序 `./pastebin`

__编译步骤：__

```sh
git clone https://github.com/nibazshab/pastebin.git
cd pastebin
CGO_ENABLED=0 go build -ldflags="-s -w"
```

## api

### POST /

参数:

 - `t`: 要存储的文本内容

响应:

 - 成功：返回存储了文本内容的链接
 - 失败：什么都不返回

### GET /{uid}

响应:

 - 成功：返回该链接所对应的文本内容
 - 失败：什么都不返回

## 计划

- [ ] 响应 `dmesg | curl -F "t=@-" host` 形式的请求

## 开源地址

https://github.com/nibazshab/pastebin

## 使用许可

MIT © ZShab Niba
