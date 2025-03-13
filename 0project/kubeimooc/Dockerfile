# 使用 golang:alpine 作为基础镜像，并将此阶段命名为 builder
# 这个阶段主要用于编译 Go 代码
FROM golang:alpine as builder

# 设置工作目录
WORKDIR /go/src/kubeimooc.com/server

# 修正：COPY 命令需要指定源和目标，这里将当前目录下的所有文件复制到工作目录
COPY . .

# 设置 GO11MODULE 环境变量为 on，强制启用 Go Modules 进行依赖管理
RUN go env -w GO11MODULE=on \
    # 设置 GOPROXY 环境变量，使用国内的代理服务器 https://goproxy.cn 加速依赖包下载，若代理无所需包则直接从源仓库下载
    && go env -w GOPROXY=https://goproxy.cn,direct \
    # 设置 CGO_ENABLE 环境变量为 0，禁用 CGO 以构建跨平台可移植的可执行文件
    && go env -w CGO_ENABLE=0 \
    # 打印当前 Go 环境变量的设置，方便确认前面设置的环境变量是否生效
    && go env \
    # 整理项目的依赖关系，根据代码导入情况添加缺失依赖、移除无用依赖，并更新 go.mod 和 go.sum 文件
    && go mod tidy \
    # 编译当前目录下的 Go 代码，使用 -o 选项指定生成的可执行文件名为 server
    && go build -o server .

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