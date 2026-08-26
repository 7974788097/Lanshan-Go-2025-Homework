# 构建阶段
FROM golang:1.25.4-alpine AS builder

# 设置工作目录
WORKDIR /build

# 安装必要的构建工具
RUN apk add --no-cache git

# 设置Go代理
ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on

# 复制go mod文件
COPY go.mod go.sum ./

# 先尝试下载，如果失败则使用本地缓存（如果有）
RUN go mod download || echo "download failed, will use local cache if available"

# 复制源代码
COPY . .

# 构建应用（使用-mod=mod参数）
RUN go build -mod=mod -o a_pile_of_shit main.go

# 运行阶段
FROM alpine:latest

# 安装时区数据
RUN apk add --no-cache tzdata

ENV TZ=Asia/Shanghai
WORKDIR /app

COPY --from=builder /build/a_pile_of_shit .
RUN mkdir -p /app/logs

EXPOSE 8080

CMD ["./a_pile_of_shit"]