
# 项目介绍 

从项目: `https://github.com/wangbjun/wake_up` 拷贝过来，合成一个文件

小米路由器3增强版使用opkg一直无法下载wol package，只能够下载`etherwake`, 但是使用 `etherwake` 无法唤醒主机，使用tcpdump发现没有发送包。

编译命令如下:

```
export GOOS=linux
export GOARCH=mipsle
go build -ldflags "-s -w" -o wol main.go   # 得到可以在小米路由器3增强版上执行的`wol`
```


