# TradingAgents Dockerfile
# 多阶段构建优化镜像大小

# ================================
# Stage 1: Build
# ================================
FROM golang:1.21-alpine AS builder

# 安装构建依赖
RUN apk add --no-cache git ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod go.sum* ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}" \
    -o /app/tradingagents \
    ./cmd/main.go

# ================================
# Stage 2: Runtime
# ================================
FROM alpine:3.19

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 创建非 root 用户
RUN addgroup -g 1000 appgroup && \
    adduser -u 1000 -G appgroup -D appuser

# 设置时区
ENV TZ=UTC

# 创建应用目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/tradingagents /app/tradingagents

# 创建数据目录
RUN mkdir -p /app/data /app/results /app/logs && \
    chown -R appuser:appgroup /app

# 切换到非 root 用户
USER appuser

# 设置环境变量
ENV APP_ENV=production \
    LOG_LEVEL=info \
    DATA_DIR=/app/data \
    RESULTS_DIR=/app/results

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ["/app/tradingagents", "version"]

# 暴露端口（如果有 API 服务）
EXPOSE 8080

# 入口点
ENTRYPOINT ["/app/tradingagents"]

# 默认命令
CMD ["--help"]
