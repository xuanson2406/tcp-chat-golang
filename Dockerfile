FROM golang:1.16-alpine

RUN mkdir /app

ADD . /app

WORKDIR /app

ENV IP="172.17.0.2:3000"

RUN go build -o client client.go
