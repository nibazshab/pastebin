# PasteBin

支持命令行使用的一个匿名数据分享平台，分享例如一段话，一篇文章，一张图片，一个压缩包，然后将返回的链接发给别人

最大限制 100MB

### 快速上手

下载 Releases 中的文件直接运行即可，默认监听 10002 端口，数据默认位于 pastebin_data 目录中

命令行可以接收的参数

参数|默认值|描述
-|-|-
-port|10002|端口号
-dir|pastebin_data|数据目录

配合 systemd 使用 pastebin.service

```ini
[Unit]
Description=pastebin service
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

### API

- ___POST /___，返回存储了数据内容的链接
  - 表单 f = 文件，上传的文件  
  - headers
    - type: 1/2/3/4，类型（可选）

- ___GET /{uid}___，返回该链接所对应的内容

- ___DELETE /{uid}___，删除记录
  - headers
    - token: xxx，令牌

```sh
cat /etc/hosts | curl -F f=@- 127.0.0.1:10002

curl 127.0.0.1:10002/fgs3

curl -X DELETE 127.0.0.1:10002/fgs3 -H 'token: 2A9B3F692B1715A6'
```
## 许可证

MIT © ZShab Niba
