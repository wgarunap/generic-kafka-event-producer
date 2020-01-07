package main

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/gorilla/mux"
	"github.com/pickme-go/log/v2"
	schemaregistry "github.com/pickme-go/schema-registry"
	"io/ioutil"
	"net/http"
	"time"
)

type subjects []string

var schemas = map[string][]int{}

var Registry *schemaregistry.Registry
var Producer sarama.SyncProducer

func main() {
	schemaurl := `http://35.208.110.113:8081`

	var err error
	Registry, err = schemaregistry.NewRegistry(schemaurl)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Get(schemaurl + `/subjects`)
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
		resp, err := http.Get(schemaurl + `/subjects/` + subject + `/versions`)
		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(body, &versions)
		if err != nil {
			log.Fatal(err)
		}

		for _, v := range versions {

			schemas[subject] = append(schemas[subject], v)

			if err := Registry.Register(subject, v, func(data []byte) (v interface{}, err error) {
				return nil, nil
			}); err != nil {
				log.Fatal(err)
			}
			log.Info(fmt.Sprintf("Registered subject:%s, version:%d", subject, v))
		}
	}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	Producer, err = sarama.NewSyncProducer([]string{`localhost:9092`}, config)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/publish", publish).Methods(http.MethodPost)

	log.Info("serving now... localhost:8000")
	http.ListenAndServe(`:8000`, router)

}

type Event struct {
	Subject string      `json:"subject"`
	Version int         `json:"version"`
	Topic   string      `json:"topic"`
	Value   interface{} `json:"value"`
	Key     string      `json:"key"`
}

func publish(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	log.Debug("Request:", string(body))

	var e = Event{}
	err = json.Unmarshal(body, &e)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	_, ok := schemas[e.Subject]
	if !ok {
		w.WriteHeader(400)
		w.Write([]byte(`subject not registered`))
		return
	}

	available := false
	for _, version := range schemas[e.Subject] {
		if version == e.Version {
			available = true
			break
		}
	}
	if !available {
		w.WriteHeader(400)
		w.Write([]byte(`requested verion not registered in the application`))
		return
	}

	encodedValue, err := Registry.WithSchema(e.Subject, e.Version).Encode(e.Value)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	p, o, err := Producer.SendMessage(&sarama.ProducerMessage{
		Topic:     e.Topic,
		Key:       sarama.ByteEncoder([]byte(e.Key)),
		Value:     sarama.ByteEncoder(encodedValue),
		Timestamp: time.Now(),
	})


	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf(`{"partition":%d, "offset":%d}`, p, o)))

}
