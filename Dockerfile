FROM golang:1.16-alpine3.14 AS builder

COPY . /opt/
WORKDIR /opt/

RUN env GO111MODULE=on GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o generic-kafka-event-producer *.go

FROM alpine:3.14.0

RUN apk update && apk add --no-cache \
    ca-certificates \
    libc6-compat \
    tzdata

ENV TZ=Asia/Colombo

WORKDIR /opt

COPY --from=builder /opt/generic-kafka-event-producer /opt

ENTRYPOINT ["sh", "-c","/opt/generic-kafka-event-producer"]
