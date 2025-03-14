# 使用 golang:alpine 作为基础镜像，并将此阶段命名为 builder
# 这个阶段主要用于编译 Go 代码
FROM golang:1.23.6 as builder

# 设置工作目录
WORKDIR /go/src/kubeimooc.com/server

# COPY 命令需要指定源和目标，这里将当前目录下的所有文件复制到工作目录
COPY . .

# 设置环境变量 GOPROXY
ENV GOPROXY=https://goproxy.cn,direct
# 下载go依赖
RUN go mod download
# 将 . 目录下代码编译为 aaa 执行文件
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# 使用 alpine:latest 作为基础镜像，这个阶段用于创建最终的运行镜像
FROM alpine:latest

# 添加镜像维护者信息
LABEL MAINTAINER="shisan"

# 设置工作目录
WORKDIR /go/src/kubeimooc.com/server

# 从 builder 阶段复制配置文件到当前工作目录
COPY --from=builder /go/src/kubeimooc.com/server/config.yaml ./config.yaml

# 从 builder 阶段复制 Kubernetes 配置文件到当前工作目录下的 .kube 目录
COPY --from=builder /go/src/kubeimooc.com/server/.kube/config ./.kube/config

# 从 builder 阶段复制编译好的可执行文件到当前工作目录
COPY --from=builder /go/src/kubeimooc.com/server/server ./

# 声明容器运行时监听的端口
EXPOSE 8082

# 设置容器启动时执行的命令
ENTRYPOINT ["./server"]