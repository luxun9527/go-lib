FROM golang:1.20.7-alpine3.18 AS builder
#https://bbs.huaweicloud.com/blogs/251201
#https://juejin.cn/post/7159939264510328840
# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 移动到工作目录：/build
WORKDIR /build

# 将代码复制到容器中
COPY . .

RUN go mod tidy

COPY . .

# 将我们的代码编译成二进制可执行文件 app
RUN go build -o bubble .

###################
# 接下来创建一个小镜像
###################
FROM scratch

COPY ./conf /conf

COPY ./wait-for.sh /

# 从builder镜像中把/ dist/app 拷贝到当前目录
COPY --from=builder /build/app /

RUN chmod 755 wait-for.sh
# 需要运行的命令
#ENTRYPOINT ["/app"]