package errors

import (
	"context"
	"net/http"
	"strconv"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/tryfix/log"
	"github.com/tryfix/metrics"
)

/* http error codes*/
const (
	//Domain
	ErrDecodingRequest = iota + 10000
	ErrValidationInRequest
	ErrAdapter
	ErrAdaptorURL
	ErrAdaptorRequestPayload
	ErrDecodingAdaptorResponse

	//Application
	ErrApplication = iota + 20000

	ErrUnknown = iota + 9999
)

// DomainError is the type of errors thrown by business logic.
type domainError struct {
	msg     string
	code    int
	details string
}

// Error returns the DomainError message.
func (e domainError) Error() string {
	return `{` +
		`"msg":"` + e.msg + `",` +
		`"code":` + strconv.Itoa(e.code) + `,` +
		`"details":"` + e.details + `"` +
		`}`
}

// NewDomainError creates a new DomainError instance.
func NewDomainError(message string, code int, details string) error {

	err := domainError{}

	err.msg = message
	err.code = code
	err.details = details

	return &err
}

// ApplicationError is
type applicationError struct {
	msg     string
	code    int
	details string
}

// ApplicationError returns the DomainError message.
func (e applicationError) Error() string {
	return `{` +
		`"msg":"` + e.msg + `",` +
		`"code":` + strconv.Itoa(e.code) + `,` +
		`"details":"` + e.details + `"` +
		`}`
}

// NewAplicationError creates a new ApplicationError instance.
func NewAplicationError(message string, code int, details string) error {
	err := applicationError{}

	err.msg = message
	err.code = code
	err.details = details

	return &err
}

// LogErrorHandler v
type LogErrorHandler struct {
	logger  log.Logger
	metrics metrics.Counter
}

// NewLogErrorHandler resurn an instance of a LogErrorHandler
func NewLogErrorHandler(logger log.Logger, reporter metrics.Reporter) *LogErrorHandler {
	return &LogErrorHandler{
		logger:  logger,
		metrics: reporter.Counter(metrics.MetricConf{Path: "endpoint_err", Labels: []string{"err"}}),
	}
}

func (h *LogErrorHandler) Handle(ctx context.Context, err error) {
	h.metrics.Count(1, map[string]string{"err": err.Error()})
	h.logger.ErrorContext(ctx, err)
}

// CustomErrEncoder encoder
func CustomErrEncoder() httptransport.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		switch err.(type) {
		case *domainError:
			w.WriteHeader(http.StatusBadRequest)
		case *applicationError:
			w.WriteHeader(http.StatusInternalServerError)
		}

		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.ErrorContext(ctx, err)
		}
	}
}
