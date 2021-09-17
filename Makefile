.PHONY: build

 GO111MODULE=on
 GOOS=linux
 GOARCH=amd64

build:
	@env GO111MODULE=${GO111MODULE} GOOS=${GOOS} GOARCH=${GOARCH} -ldflags="-s -w" go build -o generic-kafka-event-producer *.go
