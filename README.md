
# 项目介绍 

从项目: `https://github.com/wangbjun/wake_up` 拷贝过来，合成一个文件

小米路由器3增强版使用opkg一直无法下载wol package，只能够下载`etherwake`, 但是使用 `etherwake` 无法唤醒主机，使用tcpdump发现没有发送包。

编译命令如下:

```
export GOOS=linux
export GOARCH=mipsle
go build -ldflags "-s -w" -o wol main.go   # 得到可以在小米路由器3增强版上执行的`wol`
```

由于是UDP的广播包，因此存在丢包的可能，可以执行多次提高唤醒机器的概率(观察发现`mac os`上的wol就是发了多个UDP广播包)，另外广播地址应为: `192.168.31.255`

小米路由器上执行命令如: `./wol 192.168.31.255 70:85:C2:76:63:7a` 


