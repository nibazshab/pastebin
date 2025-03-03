# PasteBin

一个工具，可以用来生成纯文本网页分享文字内容，或者也可以上传附件，作为一个临时的文件服务器分享给别人

### 快速上手

直接运行独立的二进制文件即可，默认监听 10002 端口，数据默认位于 pastebin_data 目录中，单次分享的限制为 100mb

```sh
./pastebin
```

配合 systemd 使用 pastebin.service

```ini
[Unit]
Description=Pastebin service
[Service]
ExecStart=/usr/local/pastebin/pastebin
Restart=on-failure
[Install]
WantedBy=multi-user.target
```

### 构建说明

```sh
make
```

### 使用说明

命令行可以接收的参数

参数|默认值|描述
-|-|-
-port|10002|程序监听的端口号
-dir|pastebin_data|数据目录（相对程序文件的路径）

### API

- ___POST /___

请求：multipart/form-data，返回存储了数据内容的链接

body：`f=@文件` 上传的文件
headers:（选）`X-V=1` 文件可在浏览器预览

- ___GET /{uid}___

返回该链接所对应的内容

## 许可证

MIT © ZShab Niba
