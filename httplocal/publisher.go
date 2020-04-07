package httplocal

import (
	"context"
	"encoding/json"
	"fmt"
	"generic-kafka-event-producer/domain"
	"generic-kafka-event-producer/errors"
	"generic-kafka-event-producer/schemareg"
	"io/ioutil"
	"net/http"
	"time"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/vmihailenco/msgpack/v4"

	"github.com/go-kit/kit/endpoint"
	"github.com/tryfix/schemaregistry"

	"github.com/Shopify/sarama"
	"github.com/tryfix/log"
)

func DecodeRequest() httptransport.DecodeRequestFunc {
	return func(i context.Context, req *http.Request) (request interface{}, err error) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, errors.NewDomainError(err.Error(), errors.ErrDecodingRequest, "")
		}

		log.Debug("Request:", string(body))

		var e = domain.Event{}
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
		e := request.(domain.Event)

		if e.Key.Body == nil && e.Value.Body == nil {
			return nil, errors.NewAplicationError("key.body and value.body both values are null", errors.ErrValidationInRequest, "kafka message key and value both cannot be null")
		}

		key, err := MessageEncoder(e.Key)
		if err != nil {
			return nil, err
		}
		value, err := MessageEncoder(e.Value)
		if err != nil {
			return nil, err
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
			Key:       sarama.ByteEncoder(key),
			Value:     sarama.ByteEncoder(value),
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

// MessageEncoder encodes the message based on the given format
func MessageEncoder(msg domain.Message) (encodedValue []byte, err error) {
	if msg == *new(domain.Message) {
		return nil, nil
	}

	switch msg.Format {
	case `byte`, `bytes`:
		var ok bool
		value, ok := msg.Body.(string)
		if !ok {
			return nil, errors.NewAplicationError("if the format of the kafka message is byte or null, body should be string", errors.ErrApplication, "")
		}
		encodedValue = []byte(value)

	case `string`, ``:
		var ok bool
		value, ok := msg.Body.(string)
		if !ok {
			return nil, errors.NewAplicationError("if the format of the kafka message is byte or null, body should be string", errors.ErrApplication, "")
		}
		encodedValue = []byte(value)

	case `json`:
		encodedValue, err = json.Marshal(msg.Body)
		if err != nil {
			return nil, errors.NewAplicationError(err.Error(), errors.ErrApplication, "")
		}

	case `avro`:
		if msg.Schema == nil {
			available, err := schemareg.IsSchemaAvailable(msg.Subject, msg.Version)
			if !available {
				return nil, errors.NewAplicationError(err.Error(), errors.ErrApplication, "requested subject not available")
			}

			encodedValue, err = schemareg.GetRegistry().WithSchema(msg.Subject, msg.Version).Encode(msg.Body)
			if err != nil {
				return nil, errors.NewAplicationError(err.Error(), errors.ErrApplication, "error encoding the msg into avro")
			}
			return encodedValue, err
		}
		panic(`avro encoding with given schema is not implemented`)

	case `messagepack`:
		encodedValue, err = msgpack.Marshal(msg.Body)
		if err != nil {
			return nil, errors.NewAplicationError(err.Error(), errors.ErrApplication, "error encoding the msg into messagepack")
		}
	case `protobuf`:
		panic(`protobuf encoding is not implemented`)
	}
	return encodedValue, nil
}
