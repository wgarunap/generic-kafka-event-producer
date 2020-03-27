#!/bin/bash

go build -o generic-event-producer

export SCHEMAREG_URL=http://localhost:8081
export KAFKA_BROKERS=localhost:9092

./generic-event-producer