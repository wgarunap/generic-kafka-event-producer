package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/gorilla/mux"
	"github.com/pickme-go/log/v2"
	schemaregistry "github.com/pickme-go/schema-registry"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	FormatAvro = "avro"
	FormatJson = "json"
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
	config.Version = sarama.V2_2_0_0

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
	Subject string            `json:"subject"`
	Version int               `json:"version"`
	Topic   string            `json:"topic"`
	Key     string            `json:"key"`
	Headers map[string]string `json:"headers"`
	Value   interface{}       `json:"value"`
	Format  string            `json:"format"`
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

	var encodedValue []byte

	if e.Format == FormatJson {
		encodedValue, err = json.Marshal(e.Value)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}
	} else if e.Format == FormatAvro {
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
			w.Write([]byte(`requested version not registered in the application`))
			return
		}

		encodedValue, err = Registry.WithSchema(e.Subject, e.Version).Encode(e.Value)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}
	} else {
		w.WriteHeader(400)
		w.Write([]byte(errors.New("type must be avro or json").Error()))
		return
	}

	//create headers
	var headers []sarama.RecordHeader
	for s, s2 := range e.Headers {
		headers = append(headers, sarama.RecordHeader{
			Key:   []byte(s),
			Value: []byte(s2),
		})
	}

	p, o, err := Producer.SendMessage(&sarama.ProducerMessage{
		Topic:     e.Topic,
		Key:       sarama.ByteEncoder([]byte(e.Key)),
		Value:     sarama.ByteEncoder(encodedValue),
		Headers:   headers,
		Timestamp: time.Now(),
	})

	if err != nil {
		log.Error(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`{"msg":%s, "err":%s}`, `unable to produce message`, err.Error())))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf(`{"partition":%d, "offset":%d}`, p, o)))
}
