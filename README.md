# pscan

proxy scan，扫描子网内的常见局域网代理端口。

扫描端口配置可在 `scan/settings.go` 中修改

## 使用方法

安装

```bash
go install  github.com/hatrd/pscan@latest
```

运行 `pscan` 然后输入子网编号。扫描完成后会在运行目录生成一个 txt 文件。

需要更详细的设置，使用 `pscan -h` 获取帮助。可以禁用 icmp 存活检测（有的主机不回 icmp 报文）、禁用端口检测（只用 icmp 探测内网存活情况）、设置并发数量、超时时间等。
