# Стадия сборки
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0

RUN go build -o main cmd/app/main.go

# Финальный образ
FROM alpine:3.18

# Устанавливаем необходимые зависимости
RUN apk add --no-cache \
    ca-certificates \
    ffmpeg \
    python3 \
    py3-pip

# Устанавливаем yt-dlp
RUN pip3 install --upgrade yt-dlp

WORKDIR /

COPY --from=builder /app/main .

RUN chmod +x /main

CMD ["/main"]
