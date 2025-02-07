# 使用官方的 Go 镜像作为构建环境
FROM golang:1.22 AS builder

# 设置工作目录
WORKDIR /app

# 复制 Go 模块文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源码并构建
COPY . .
# 将 . 目录下代码编译为 aaa 执行文件
RUN CGO_ENABLED=0 GOOS=linux go build -o aaa .

# 使用轻量级的 Alpine 镜像作为运行环境
FROM alpine:latest

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制可执行文件
COPY --from=builder /app/aaa .

# 暴露端口（根据你的程序调整）
EXPOSE 8080

# 运行程序
CMD ["./aaa"]


