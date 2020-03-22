package schemareg

import (
	"encoding/json"
	"fmt"
	"generic-kafka-event-producer/config"
	"io/ioutil"
	"net/http"

	"github.com/tryfix/log"
	"github.com/tryfix/schemaregistry"
)

type subjects []string

var schemas = map[string]struct{}{}

var reg *schemaregistry.Registry

func GetRegistry() *schemaregistry.Registry {
	return reg
}

func Init() {
	r, err := schemaregistry.NewRegistry(config.Config.SchemaRegUrl)
	if err != nil {
		log.Fatal(err)
	}
	reg = r
}

// RegisterEvents is
func RegisterEvents() {
	var err error

	resp, err := http.Get(config.Config.SchemaRegUrl + `/subjects`)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	sub := subjects{}
	err = json.Unmarshal(body, &sub)
	if err != nil {
		log.Fatal(err)
	}

	for _, subject := range sub {
		schemas[subject] = struct{}{}

		err := RegisterEvent(subject, int(schemaregistry.VersionAll))
		if err != nil {
			log.Fatal(err)
		}
		log.Info(fmt.Sprintf("Registered subject:%s", subject))
	}
}

func RegisterEvent(subject string, version int) error {
	err := reg.Register(subject, version,
		func(data []byte) (v interface{}, err error) {
			return nil, nil
		})
	return err
}

func IsSchemaAvailable(subject string, version int) (bool, error) {
	_, ok := schemas[subject]
	if !ok {
		err := RegisterEvent(subject, int(schemaregistry.VersionAll))
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return true, nil
}
