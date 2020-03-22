package config

import (
	"github.com/caarlos0/env"
	"github.com/tryfix/log"
)

type Conf struct {
	SchemaRegUrl string   `env:"SCHEMAREG_URL"`
	Brokers      []string `env:"KAFKA_BROKERS"`
	Port         int      `env:"PORT" envDefault:"8000"`
}

var Config Conf

func (*Conf) Register() {
	err := env.Parse(&Config)
	if err != nil {
		log.Fatal("error loading schema config, ", err)
	}
}

func (*Conf) Validate() {
	if Config.SchemaRegUrl == "" {
		log.Fatal("schema registry url configuration is not found, please set SCHEMAREG_URL to schema registry URL")
	}
	if Config.Brokers == nil {
		log.Fatal("KAFKA_BROKERS environment variable not set in the configuration")
	}
}

func (*Conf) Print() interface{} {
	defer log.Info("configs loaded")
	return Config
}
