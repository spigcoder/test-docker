# ==========================================
# 阶段 1: Builder (编译阶段)
# ==========================================
FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk --no-cache add tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o argo-server \
    ./argo

# ==========================================
# 阶段 2: Runner (运行阶段)
# ==========================================
# 使用 distroless static，不含 Shell，体积极小，最安全
FROM gcr.io/distroless/static-debian12

# 从 Builder 阶段拷贝时区数据 (解决 loc=Local 问题)
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Asia/Shanghai

# 从 Builder 阶段拷贝编译好的二进制文件
COPY --from=builder /app/argo-server /app/argo-server

# 最佳实践：使用非 root 用户运行 (distroless 内置了 nonroot 用户: uid 65532)
USER nonroot:nonroot

# 声明端口 (仅作文档说明，实际端口由代码决定)
EXPOSE 8080

# 启动命令
ENTRYPOINT ["/app/argo-server"]