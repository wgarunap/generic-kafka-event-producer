package httplocal

import (
	"context"
	"generic-kafka-event-producer/config"
	"generic-kafka-event-producer/errors"
	"generic-kafka-event-producer/httplocal/schemas"
	"generic-kafka-event-producer/httplocal/topics"
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
	router.Handle("/publish", avroHandler()).Methods(http.MethodPost)
	router.Handle("/topics", topicListHandler()).Methods(http.MethodGet)
	router.Handle("/avrosubjects", schemasListHandler()).Methods(http.MethodGet)

	log.Info("serving now... localhost:" + strconv.Itoa(config.Config.Port))

	http.ListenAndServe(`:`+strconv.Itoa(config.Config.Port), router)
}

func avroHandler() http.Handler {
	return httptransport.NewServer(
		Endpoint(
			schemareg.GetRegistry(),
			producer.GetProducer(),
		),
		DecodeRequest(),
		EncodeResponse(),
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

func topicListHandler() http.Handler {
	return httptransport.NewServer(
		topics.Endpoint(),
		topics.DecodeRequest(),
		topics.EncodeResponse(),
		httptransport.ServerErrorHandler(
			errors.NewLogErrorHandler(
				log.StdLogger,
				metrics.NoopReporter(),
			)),
		httptransport.ServerErrorEncoder(
			errors.CustomErrEncoder(),
		),
	)
}

func schemasListHandler() http.Handler {
	return httptransport.NewServer(
		schemas.Endpoint(),
		schemas.DecodeRequest(),
		schemas.EncodeResponse(),
		httptransport.ServerErrorHandler(
			errors.NewLogErrorHandler(
				log.StdLogger,
				metrics.NoopReporter(),
			)),
		httptransport.ServerErrorEncoder(
			errors.CustomErrEncoder(),
		),
	)
}

func TraceableContext() httptransport.RequestFunc {
	return func(i context.Context, request *http.Request) context.Context {
		u, _ := uuid.NewUUID()
		return traceable_context.FromContextWithUUID(i, u)
	}
}
