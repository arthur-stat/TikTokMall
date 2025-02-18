FROM golang:1.19-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o auth_service

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/auth_service .
COPY --from=builder /app/conf/dev/conf.yaml ./conf/dev/

EXPOSE 8080

CMD ["./auth_service"] 