package main

import (
	"add-service/endpoint"
	"add-service/service"
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var ints, chars metrics.Counter
	{
		// Business-level metrics.
		ints = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "go_kit",
			Subsystem: "add_service",
			Name:      "integers_summed",
			Help:      "Total count of integers summed via the Sum method.",
		}, []string{})

		chars = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "go_kit",
			Subsystem: "add_service",
			Name:      "characters_concatenated",
			Help:      "Total count of characters concatenated via the Concat method.",
		}, []string{})
	}

	var duration metrics.Histogram
	{
		// Endpoint-level metrics.
		duration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "go_kit",
			Subsystem: "add_service",
			Name:      "request_duration_seconds",
			Help:      "Request duration in seconds.",
		}, []string{"method", "success"})
	}

	svc := service.NewBasicService()
	svc = service.LoggingMiddleware(logger)(svc)
	svc = service.InstrumentingMiddleware(ints, chars)(svc)

	endpoints := endpoint.New(svc, logger, duration)

	sumHandler := httptransport.NewServer(endpoints.SumEndpoint, decodeHTTPSumRequest, encodeHTTPGenericResponse)
	concatHandler := httptransport.NewServer(endpoints.ConcatEndpoint, decodeHTTPConcatRequest, encodeHTTPGenericResponse)

	http.Handle("/sum", sumHandler)
	http.Handle("/concat", concatHandler)
	// ../prometheus/prometheus --config.file=.config/prometheus.yml
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}

func decodeHTTPSumRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoint.SumRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	return request, err
}

func decodeHTTPConcatRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoint.ConcatRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	return request, err
}

func encodeHTTPGenericResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
