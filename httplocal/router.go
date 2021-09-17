package httplocal

import (
	"context"
	"fmt"
	"generic-kafka-event-producer/config"
	"generic-kafka-event-producer/errors"
	"generic-kafka-event-producer/httplocal/schemas"
	"generic-kafka-event-producer/httplocal/topics"
	"generic-kafka-event-producer/producer"
	"generic-kafka-event-producer/schemareg"
	"net/http"
	"os"
	"os/signal"
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

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Kill, os.Interrupt)

	srv := http.Server{
		Addr:    `:` + strconv.Itoa(config.Config.Port),
		Handler: router,
	}

	go func() {
		log.Info(`server is starting on port:(` + srv.Addr + ")")
		err := srv.ListenAndServe()
		if err != nil {
			if err != http.ErrServerClosed {
				log.Error(fmt.Sprintf("error closing the server, err:%v", err))
			}
		}
	}()

	<-sig
	err := srv.Close()
	if err != nil {
		log.Error(fmt.Sprintf("error closing the server, err:%v", err))
	}

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
