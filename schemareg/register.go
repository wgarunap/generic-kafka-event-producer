package schemareg

import (
	"encoding/json"
	"fmt"
	"generic-kafka-event-producer/config"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/tryfix/log"
	"github.com/tryfix/schemaregistry"
)

type subjects []string

var schemas = map[string]map[int]struct{}{}
var mu sync.RWMutex

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

		var versions []int
		versions, err := getAllSchemaVersions(subject)
		if err != nil {
			log.Fatal(err)
		}

		for _, version := range versions {
			err := registerEvent(subject, version)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func registerEvent(subject string, version int) error {
	if err := reg.Register(subject, version, func(data []byte) (v interface{}, err error) {
		return nil, nil
	}); err != nil {
		return err
	}

	mu.Lock()
	_, ok := schemas[subject]
	if !ok {
		schemas[subject] = map[int]struct{}{}
	}
	schemas[subject][version] = struct{}{}
	mu.Unlock()

	log.Info(fmt.Sprintf("schema registered, subject:%s, version:%d", subject, version))

	return nil
}

func getAllSchemaVersions(subject string) (versions []int, err error) {
	resp, err := http.Get(config.Config.SchemaRegUrl + `/subjects/` + subject + `/versions`)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &versions)
	if err != nil {
		return nil, err
	}
	return versions, nil
}

func IsSchemaAvailable(subject string, version int) (bool, error) {
	mu.RLock()
	sub, ok := schemas[subject]
	mu.RUnlock()
	if !ok {
		var versions []int
		versions, err := getAllSchemaVersions(subject)
		if err != nil {
			return false, err
		}

		for _, version := range versions {
			err := registerEvent(subject, version)
			if err != nil {
				return false, err
			}
		}
	}

	_, ok = sub[version]
	if !ok {
		err := registerEvent(subject, version)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}
