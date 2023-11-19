package main

import (
	"add-service/service"
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
		// ts=2023-11-19T06:19:04.945776Z method=Concat a=leo b=ken v=leoken err=null
		logger = log.With(logger, "caller", log.DefaultCaller)
		// ts=2023-11-19T06:20:31.612571Z caller=middleware.go:32 method=Concat a=leo b=ken v=leoken err=null
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

	svc := service.NewBasicService()
	svc = service.LoggingMiddleware(logger)(svc)
	svc = service.InstrumentingMiddleware(ints, chars)(svc)

	sumHandler := httptransport.NewServer(MakeSumEndpoint(svc), decodeHTTPSumRequest, encodeHTTPGenericResponse)
	concatHandler := httptransport.NewServer(MakeConcatEndpoint(svc), decodeHTTPConcatRequest, encodeHTTPGenericResponse)

	http.Handle("/sum", sumHandler)
	http.Handle("/concat", concatHandler)
	// ../prometheus/prometheus --config.file=.config/prometheus.yml
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
