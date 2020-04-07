package main

import (
	"generic-kafka-event-producer/config"
	"generic-kafka-event-producer/httplocal"
	"generic-kafka-event-producer/producer"
	"generic-kafka-event-producer/schemareg"

	"github.com/wgarunap/goconf"
)

func main() {
	goconf.Load(
		new(config.Conf),
	)

	if config.Config.SchemaRegUrl != "" {
		schemareg.Init()
		schemareg.RegisterEvents()
	}

	producer.InitProducer()

	httplocal.Start()

}
