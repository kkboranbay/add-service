package endpoint

import (
	"add-service/service"
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/log"
)

// Set collects all of the endpoints that compose an add service. It's meant to
// be used as a helper struct, to collect all of the endpoints into a single
// parameter.
type Set struct {
	SumEndpoint    endpoint.Endpoint
	ConcatEndpoint endpoint.Endpoint
}

func New(svc service.Service, logger log.Logger, duration metrics.Histogram) Set {
	var sumEndpoint endpoint.Endpoint
	{
		sumEndpoint = MakeSumEndpoint(svc)
		sumEndpoint = LoggingMiddleware(log.With(logger, "method", "Sum"))(sumEndpoint)
		sumEndpoint = InstrumentingMiddleware(duration.With("method", "Sum"))(sumEndpoint)
	}

	var concatEndpoint endpoint.Endpoint
	{
		concatEndpoint = MakeConcatEndpoint(svc)
		concatEndpoint = LoggingMiddleware(log.With(logger, "method", "Concat"))(concatEndpoint)
		concatEndpoint = InstrumentingMiddleware(duration.With("method", "Concat"))(concatEndpoint)
	}

	return Set{
		SumEndpoint:    sumEndpoint,
		ConcatEndpoint: concatEndpoint,
	}
}

// MakeSumEndpoint constructs a Sum endpoint wrapping the service.
func MakeSumEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SumRequest)
		v, err := s.Sum(ctx, req.A, req.B)
		return SumResponse{V: v, Err: err}, nil
	}
}

// MakeConcatEndpoint constructs a Concat endpoint wrapping the service.
func MakeConcatEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ConcatRequest)
		v, err := s.Concat(ctx, req.A, req.B)
		return ConcatResponse{V: v, Err: err}, nil
	}
}

type SumRequest struct {
	A, B int
}

type SumResponse struct {
	V   int   `json:"v"`
	Err error `json:"-"`
}

type ConcatRequest struct {
	A, B string
}

type ConcatResponse struct {
	V   string `json:"v"`
	Err error  `json:"-"`
}
