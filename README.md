# PasteBin

支持命令行使用的一个匿名数据分享平台，在这里分享例如一段话，一篇文章，一张图片，一个压缩包，然后将返回的链接发给别人即可

最大限制 100MB

### 快速上手

直接运行即可，默认监听 10002 端口，数据默认位于 pastebin_data 目录中

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

body `f=@文件` 上传的文件  
headers `X-V=1` 文件可在浏览器预览（选）

- ___GET /{uid}___

返回该链接所对应的内容

## 许可证

MIT © ZShab Niba
