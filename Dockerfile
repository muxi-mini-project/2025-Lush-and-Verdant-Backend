# 使用官方 Go 镜像作为构建环境
FROM golang:1.20 AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 并下载依赖
COPY go.mod go.sum ./
RUN go mod tidy

# 复制所有源代码到容器
COPY . .

# 编译 Go 应用，生成可执行文件
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 使用一个轻量级镜像运行应用
FROM alpine:latest

# 设置工作目录
WORKDIR /root/

# 复制编译好的二进制文件
COPY --from=builder /app/main .

# 暴露端口（与 Go 代码中的端口一致）
EXPOSE 8080

# 运行应用
CMD ["./main"]