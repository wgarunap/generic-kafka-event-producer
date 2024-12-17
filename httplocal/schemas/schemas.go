package schemas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"generic-kafka-event-producer/schemareg"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func DecodeRequest() httptransport.DecodeRequestFunc {
	return func(i context.Context, req *http.Request) (request interface{}, err error) {
		return nil, nil
	}
}
func Endpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		schemas := schemareg.GetAllSchemas()
		return schemas, nil
	}
}

type Response struct {
	Data struct {
		Subjects []Subject `json:"subjects"`
	} `json:"data"`
}

type Subject struct {
	Subject  string `json:"subject"`
	Versions []int  `json:"versions"`
}

func EncodeResponse() httptransport.EncodeResponseFunc {
	return func(context context.Context, writer http.ResponseWriter, response interface{}) error {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(200)

		res := Response{}

		schemas := response.(map[string][]int)
		for s, v := range schemas {
			res.Data.Subjects = append(res.Data.Subjects, Subject{
				Subject:  s,
				Versions: v,
			})
		}

		b, _ := json.Marshal(res)
		fmt.Fprintf(writer, `%s`, b)
		return nil
	}
}
