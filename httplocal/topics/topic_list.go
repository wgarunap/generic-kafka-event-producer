package topics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"generic-kafka-event-producer/config"
	"generic-kafka-event-producer/errors"

	"github.com/IBM/sarama"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/tryfix/log"
)

func DecodeRequest() httptransport.DecodeRequestFunc {
	return func(i context.Context, req *http.Request) (request interface{}, err error) {
		return nil, nil
	}
}
func Endpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		//admin connection config
		saramaConfig := sarama.NewConfig()
		saramaConfig.Version = sarama.V2_4_0_0

		//source cluster admin conncetion
		sourceClusterAdmin, err := sarama.NewClusterAdmin(config.Config.Brokers, saramaConfig)
		if err != nil {
			log.Fatal(err)
		}

		topics, err := sourceClusterAdmin.ListTopics()
		if err != nil {
			log.Fatal(err)
		}

		if err != nil {
			return nil, errors.NewAplicationError(err.Error(), errors.ErrAdapter, `unable to produce message`)
		}

		return topics, nil
	}
}

type Response struct {
	Data struct {
		Topics []ResponseTopic `json:"topics"`
	} `json:"data"`
}
type ResponseTopic struct {
	Name               string `json:"name"`
	NumberOfPartitions int    `json:"numberOfPartitions"`
}

func EncodeResponse() httptransport.EncodeResponseFunc {
	return func(context context.Context, writer http.ResponseWriter, response interface{}) error {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(200)

		res := Response{}

		topics := response.(map[string]sarama.TopicDetail)
		for t, v := range topics {
			res.Data.Topics = append(res.Data.Topics, ResponseTopic{t, int(v.NumPartitions)})
		}
		// resp, ok := response.(string)

		b, _ := json.Marshal(res)
		fmt.Fprintf(writer, `%s`, b)
		return nil
	}
}
