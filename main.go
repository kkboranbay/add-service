package main

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	svc := NewBasicService()

	sumHandler := httptransport.NewServer(MakeSumEndpoint(svc), decodeHTTPSumRequest, encodeHTTPGenericResponse)
	concatHandler := httptransport.NewServer(MakeConcatEndpoint(svc), decodeHTTPConcatRequest, encodeHTTPGenericResponse)

	http.Handle("/sum", sumHandler)
	http.Handle("/concat", concatHandler)
	http.ListenAndServe(":8080", nil)
}
