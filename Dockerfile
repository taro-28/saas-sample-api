# syntax=docker/dockerfile:1

FROM golang:1.21.5-alpine

ARG DSN
ENV DSN=${DSN}

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /server

EXPOSE 8080

CMD ["/server"]