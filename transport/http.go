package transport

import (
	"add-service/endpoint"
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHTTPHandler(endpoints endpoint.Set) http.Handler {
	m := http.NewServeMux()

	m.Handle("/sum", httptransport.NewServer(
		endpoints.SumEndpoint,
		decodeHTTPSumRequest,
		encodeHTTPGenericResponse,
	))

	m.Handle("/concat", httptransport.NewServer(
		endpoints.ConcatEndpoint,
		decodeHTTPConcatRequest,
		encodeHTTPGenericResponse,
	))

	return m
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

// encodeHTTPGenericResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func encodeHTTPGenericResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
