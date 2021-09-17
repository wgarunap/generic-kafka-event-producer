package main

import (
	"generic-kafka-event-producer/config"
	"generic-kafka-event-producer/httplocal"
	"generic-kafka-event-producer/producer"
	"generic-kafka-event-producer/schemareg"
	"github.com/tryfix/log"
	"github.com/wgarunap/goconf"
)

func main() {
	goconf.Load(
		new(config.Conf),
	)

	log.StdLogger = log.NewLog(
		log.WithLevel(config.Config.LogLevel),
		log.WithColors(config.Config.LogColor),
	).Log()

	if config.Config.SchemaRegUrl != "" {
		schemareg.Init()
		schemareg.RegisterEvents()
	}

	// initializing the kafka producer
	producer.InitProducer()

	// blocking call to serve requests
	httplocal.Start()

}
