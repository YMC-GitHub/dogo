# ==============================================
# 【构建阶段】同时编译 Alpine(Linux) + Windows 版本
# ==============================================
FROM alpine:latest AS builder

# 安装 Go 编译环境 + upx 压缩
RUN apk add --no-cache go upx

WORKDIR /app

# 复制源码
COPY ./cmd/dogo/ ./cmd/dogo/
COPY ./ipkg/dogo/ ./ipkg/dogo/

RUN [ -f go.mod ] || go mod init github.com/ymc-github/dogo
RUN go mod tidy

# ------------------------------
# 1. 编译 Linux / Alpine 平台 (amd64)
# ------------------------------
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -extldflags=-static" -trimpath -o dogo ./cmd/dogo && \
    upx --best --ultra-brute dogo

# ------------------------------
# 2. 编译 Windows 平台 (amd64)
# ------------------------------
RUN CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-w -s -extldflags=-static" -trimpath -o dogo.exe ./cmd/dogo && \
    upx --best --ultra-brute dogo.exe

# ==============================================
# 【开发环境】dev（包含两个平台的二进制）
# ==============================================
FROM alpine:latest AS dev

COPY --from=builder /app/dogo     /usr/local/bin/
COPY --from=builder /app/dogo.exe /usr/local/bin/

RUN apk add --no-cache bash curl vim tzdata openrc chrony
CMD ["sleep", "infinity"]

# ==============================================
# 【生产环境】pro（仅 Linux 版）
# ==============================================
FROM alpine:latest AS runtime


RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 直接安装 tzdata，保持安装状态（仅增加约 2MB）
RUN apk add --no-cache tzdata bash

# 设置时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY --from=builder /app/dogo /usr/local/bin/

# RUN dogo Asia/Shanghai
CMD ["sleep", "infinity"]

# ==============================================
# 【导出】空镜像，同时输出两个二进制
# ==============================================
FROM scratch AS export
COPY --from=builder /app/dogo      ./
COPY --from=builder /app/dogo.exe ./
CMD ["sleep", "infinity"]
