# syntax=docker/dockerfile:1

FROM golang:1.21.3-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

# バイナリファイルにビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o /server

EXPOSE 8080

# バイナリファイルを実行
CMD ./server