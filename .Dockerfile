# ビルド用ステージ
FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o myapp

# 実行用ステージ
FROM gcr.io/distroless/base-debian11
WORKDIR /root/
COPY --from=builder /app/myapp .
CMD ["./myapp"]
