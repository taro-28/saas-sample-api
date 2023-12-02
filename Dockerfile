# syntax=docker/dockerfile:1

FROM golang:1.21.3-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /server

EXPOSE 8080

CMD ./server