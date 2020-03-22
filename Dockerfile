FROM golang:1.14.0 as builder

COPY . /opt/
WORKDIR /opt/

RUN env GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o generic-kafka-event-producer main.go

FROM alpine:3.10

RUN apk update && apk add --no-cache \
    ca-certificates \
    libc6-compat \
    tzdata

ENV TZ=Asia/Colombo

WORKDIR /opt

COPY --from=builder /opt/generic-kafka-event-producer . 
RUN ls -la
ENTRYPOINT ["sh", "-c","./generic-kafka-event-producer"]