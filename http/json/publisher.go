package json

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"generic-kafka-event-producer/errors"

	"github.com/Shopify/sarama"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/tryfix/log"
	"github.com/tryfix/schemaregistry"
)

type Event struct {
	Subject string            `json:"subject"`
	Version int               `json:"version"`
	Topic   string            `json:"topic"`
	Key     string            `json:"key"`
	Headers map[string]string `json:"headers"`
	Value   interface{}       `json:"value"`
}

func DecodeRequest() httptransport.DecodeRequestFunc {
	return func(i context.Context, req *http.Request) (request interface{}, err error) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, errors.NewDomainError(err.Error(), errors.ErrDecodingRequest, "")
		}

		log.Debug("Request:", string(body))

		var e = Event{}
		err = json.Unmarshal(body, &e)
		if err != nil {
			return nil, errors.NewDomainError(err.Error(), errors.ErrDecodingRequest, "")
		}
		return e, nil
	}
}

func Endpoint(
	registry *schemaregistry.Registry,
	producer sarama.SyncProducer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var e = Event{}

		var encodedValue []byte

		encodedValue, err = json.Marshal(e.Value)
		if err != nil {
			return nil, errors.NewAplicationError(err.Error(), errors.ErrApplication, "")
		}
		//create headers
		var headers []sarama.RecordHeader
		for s, s2 := range e.Headers {
			headers = append(headers, sarama.RecordHeader{
				Key:   []byte(s),
				Value: []byte(s2),
			})
		}

		p, o, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic:     e.Topic,
			Key:       sarama.ByteEncoder([]byte(e.Key)),
			Value:     sarama.ByteEncoder(encodedValue),
			Headers:   headers,
			Timestamp: time.Now(),
		})

		if err != nil {
			return nil, errors.NewAplicationError(err.Error(), errors.ErrApplication, `unable to produce message`)
		}

		return fmt.Sprintf(`{"partition":%d, "offset":%d}`, p, o), nil
	}
}

func EncodeResponse() httptransport.EncodeResponseFunc {
	return func(context context.Context, writer http.ResponseWriter, response interface{}) error {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(200)
		fmt.Fprintf(writer, response.(string))
		return nil
	}
}
