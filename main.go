package main

import (
	"add-service/service"
	"net/http"
	"os"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
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

	svc := service.NewBasicService()
	svc = service.LoggingMiddleware(logger)(svc)

	sumHandler := httptransport.NewServer(MakeSumEndpoint(svc), decodeHTTPSumRequest, encodeHTTPGenericResponse)
	concatHandler := httptransport.NewServer(MakeConcatEndpoint(svc), decodeHTTPConcatRequest, encodeHTTPGenericResponse)

	http.Handle("/sum", sumHandler)
	http.Handle("/concat", concatHandler)
	http.ListenAndServe(":8080", nil)
}
