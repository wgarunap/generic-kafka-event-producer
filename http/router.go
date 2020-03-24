package http

import (
	"context"
	"generic-kafka-event-producer/config"
	"generic-kafka-event-producer/errors"
	"generic-kafka-event-producer/http/avro"
	"generic-kafka-event-producer/http/json"
	nullpublisher "generic-kafka-event-producer/http/null_publisher"
	"generic-kafka-event-producer/producer"
	"generic-kafka-event-producer/schemareg"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/gorilla/mux"
	"github.com/tryfix/log"
	"github.com/tryfix/metrics"
	traceable_context "github.com/tryfix/traceable-context"
)

// Start the http server
func Start() {
	router := mux.NewRouter()
	router.Handle("/publish/avro", avroHandler()).Methods(http.MethodPost)
	router.Handle("/publish/json", jsonHandler()).Methods(http.MethodPost)
	router.Handle("/publish/null", nullHandler()).Methods(http.MethodPost)

	log.Info("serving now... localhost:8000")

	http.ListenAndServe(`:`+strconv.Itoa(config.Config.Port), router)
}

func avroHandler() http.Handler {
	return httptransport.NewServer(
		avro.Endpoint(
			schemareg.GetRegistry(),
			producer.GetProducer(),
		),
		avro.DecodeRequest(),
		avro.EncodeResponse(),
		httptransport.ServerErrorHandler(
			errors.NewLogErrorHandler(
				log.StdLogger,
				metrics.NoopReporter(),
			)),
		httptransport.ServerErrorEncoder(
			errors.CustomErrEncoder(),
		),
		httptransport.ServerBefore(TraceableContext()),
	)
}
func jsonHandler() http.Handler {
	return httptransport.NewServer(
		json.Endpoint(
			schemareg.GetRegistry(),
			producer.GetProducer(),
		),
		json.DecodeRequest(),
		json.EncodeResponse(),
		httptransport.ServerErrorHandler(
			errors.NewLogErrorHandler(
				log.StdLogger,
				metrics.NoopReporter(),
			)),
		httptransport.ServerErrorEncoder(
			errors.CustomErrEncoder(),
		),
		httptransport.ServerBefore(TraceableContext()),
	)
}

func nullHandler() http.Handler {
	return httptransport.NewServer(
		nullpublisher.Endpoint(
			producer.GetProducer(),
		),
		nullpublisher.DecodeRequest(),
		nullpublisher.EncodeResponse(),
		httptransport.ServerErrorHandler(
			errors.NewLogErrorHandler(
				log.StdLogger,
				metrics.NoopReporter(),
			)),
		httptransport.ServerErrorEncoder(
			errors.CustomErrEncoder(),
		),
		httptransport.ServerBefore(TraceableContext()),
	)
}

func TraceableContext() httptransport.RequestFunc {
	return func(i context.Context, request *http.Request) context.Context {
		u, _ := uuid.NewUUID()
		return traceable_context.FromContextWithUUID(i, u)
	}
}
