FROM golang:1.23.4-alpine AS builder

# 设置 GOPROXY
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY . .

# 安装 git（一些依赖可能需要）
RUN apk add --no-cache git

RUN go mod download
RUN go build -o order_service

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/order_service .
COPY --from=builder /app/conf/dev/conf.yaml ./conf/dev/

EXPOSE 8080

CMD ["./order_service"] 