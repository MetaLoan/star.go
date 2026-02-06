FROM golang:1.21-alpine AS build
WORKDIR /src

# 安装 CGO 和编译依赖
RUN apk add --no-cache git gcc musl-dev make

# 创建目标目录
RUN mkdir -p /usr/local/lib /usr/local/include

# 从 GitHub 克隆并编译 Swiss Ephemeris
RUN git clone --depth 1 https://github.com/aloistr/swisseph.git /swe && \
    cd /swe && \
    gcc -c -O2 -Wall -fPIC swecl.c swedate.c swehel.c swehouse.c swejpl.c \
        swemmoon.c swemplan.c sweph.c swephlib.c && \
    ar rcs libswe.a *.o && \
    cp libswe.a /usr/local/lib/ && \
    cp swephexp.h sweph.h sweodef.h /usr/local/include/

# 复制 go mod 文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 设置 CGO 环境变量并构建
ENV CGO_ENABLED=1
ENV CGO_LDFLAGS="-L/usr/local/lib -lswe -lm"
ENV CGO_CFLAGS="-I/usr/local/include"

RUN go build -tags swe -o /out/star .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates curl
WORKDIR /app
ENV GIN_MODE=release
EXPOSE 8080
COPY --from=build /out/star /app/star
CMD ["/app/star"]
