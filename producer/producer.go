package producer

import (
	"log"

	"generic-kafka-event-producer/config"
	"github.com/IBM/sarama"
)

var producer sarama.SyncProducer

func InitProducer() {
	conf := sarama.NewConfig()
	conf.Producer.Return.Successes = true
	conf.Version = sarama.V2_2_0_0

	var err error
	producer, err = sarama.NewSyncProducer(config.Config.Brokers, conf)
	if err != nil {
		log.Fatal(err)
	}
}

func GetProducer() sarama.SyncProducer {
	return producer
}
