#!/usr/bin/env bash

binary=generic-kafka-event-producer

rm -rf ${binary}

go build -o $binary

if [[ ! -f $binary ]]; then 
    echo "build not available"
    exit 1
fi

echo "starting the application..."

export SCHEMAREG_URL=http://localhost:8081
export KAFKA_BROKERS=localhost:9092
./$binary
