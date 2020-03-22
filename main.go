package main

import (
	"generic-kafka-event-producer/config"
	"generic-kafka-event-producer/http"
	"generic-kafka-event-producer/producer"
	"generic-kafka-event-producer/schemareg"

	"github.com/wgarunap/goconf"
)

func main() {
	goconf.Load(
		new(config.Conf),
	)

	schemareg.Init()
	schemareg.RegisterEvents()

	producer.InitProducer()

	http.Start()

}
