package gapi

import (
	"log"
	"net/http"
	"time"
)

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func (rec *ResponseRecorder) Write(body []byte) (int, error) {
	rec.Body = body
	return rec.ResponseWriter.Write(body)
}

func HTTPLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		rec := &ResponseRecorder{
			ResponseWriter: res,
			StatusCode:     http.StatusOK,
		}
		handler.ServeHTTP(rec, req)
		duration := time.Since(startTime)

		log.Println(duration)
		//logger := log.Info()
		//if rec.StatusCode != http.StatusOK {
		//	logger = log.Error().Bytes("body", rec.Body)
		//}
		//
		//logger.Str("protocol", "http").
		//	Str("method", req.Method).
		//	Str("path", req.RequestURI).
		//	Int("status_code", rec.StatusCode).
		//	Str("status_text", http.StatusText(rec.StatusCode)).
		//	Dur("duration", duration).
		//	Msg("received a HTTP request")
	})
}
