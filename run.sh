#!/bin/sh

binary=generic-kafka-event-producer

go build -o $binary

if [[ ! -f $binary ]]; then 
    echo "build not available"
    exit 1
fi

echo "starting the application..."

SCHEMAREG_URL=http://localhost:8081 \
KAFKA_BROKERS=localhost:9092 \
./$binary
