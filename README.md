# helloworld
a simple web server by golang, it be made to a docker image!

这个工程我已经编译好了helloworld可执行文件，所以您在docker build的时候可以直接使用了。

如果您修改了helloworld.go文件，然后想要重新编译它，则不能在你的linux系统(比如ubuntu、centos等)中直接编译。那样编译出来的可执行文件，build进docker之后是无法运行的，docker启动会报错No such file or directory。因为我们docker中使用的是alpine环境，而不是golang环境，无法运行。

想要正确的编译出在docker下的alpine环境中可用，参考下文章

# 使用alpinelinux 构建 golang http 启动了，才15mb

1，关于alpine 环境
http://blog.csdn.net/freewebsys/article/details/53615757 
昨天研究了下golang的http服务器。 
发现在启动的时候报错：

No such file or directory

发现这个错误，开始还以为是alpine 的系统lib库少了， 
必须使用docker 官方的golang镜像呢。 
后来研究明白了，其实是因为我的宿主是centos。 
我在centos 上编译了 golang，然后拷贝到alpine 环境上造成的。 
解决办法。 
1，使用golang:alpine 镜像 241 mb 进行编译，映射一个文件夹。 
2，然后把编译好的文件拷贝出来，放到alpine的镜像上即可。

这样一个15.24 MB golang 环境就好了。 
因为还安装了一个 bash ，可以进入系统查看。

2，操作流程
首先构建一个golang build 的环境。
```
FROM       docker.io/golang:alpine

RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.4/main" > /etc/apk/repositories

RUN apk add --update curl bash && \
    rm -rf /var/cache/apk/*
```
编译镜像：
```
docker build -t demo/go-build:1.0 .
```
启动镜像，并把/data/go 目录映射到 /data/go目录，其中–rm 表示退出之后删除镜像。
```
docker run -it -v /data/gocode:/data/gocode --rm demo/go-build:1.0 /bin/bash
#cd /data/go
#go build http.go
```
其中http.go 文件：
```
package main

import (
        "fmt"
        "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
        http.HandleFunc("/", handler)
        http.ListenAndServe("0.0.0.0:8080", nil)
}
```
摘自golang 官方的httpdemo。

3，将alpine 和go http打包
在alpine环境下编译的http 包再做一个镜像，拷贝到alpine系统下：
```
FROM       docker.io/alpine:latest
MAINTAINER demo <juest a demo>

RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.4/main" > /etc/apk/repositories

RUN apk add --update curl bash && \
    rm -rf /var/cache/apk/*

RUN mkdir -p /data/go
COPY http /data/go

EXPOSE 8080

ENTRYPOINT ["/data/go/http"]
```
打包，并把http 启动。
```
docker build -t demo/go-http:1.0 .
docker run -d -p 8080:8080 --name go-http demo/go-http:1.0
```
直接访问 curl localhost:8080 即可了。


