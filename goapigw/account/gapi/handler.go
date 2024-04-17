package gapi

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
	"strconv"
)

func httpResponseModifier(ctx context.Context, w http.ResponseWriter, p proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	// set http status code
	if vals := md.HeaderMD.Get("x-http-code"); len(vals) > 0 {
		code, err := strconv.Atoi(vals[0])
		if err != nil {
			return err
		}
		// delete the headers to not expose any grpc-metadata in http response
		delete(md.HeaderMD, "x-http-code")
		delete(w.Header(), "Grpc-Metadata-X-Http-Code")
		w.WriteHeader(code)
	}

	return nil
}

func incomingHeaderMatcher(header string) (string, bool) {
	log.Println("header:", header)
	return "", true
}

func metadataMatcher(ctx context.Context, req *http.Request) metadata.MD {
	bearer := req.Header.Get("Authorization")
	log.Println("token:", bearer)
	return nil
}

func headerMatcher(header string) (string, bool) {
	log.Println("header:", header)
	return "", true
}

func routingErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, writer http.ResponseWriter, request *http.Request, statusCode int) {
	switch statusCode {
	case http.StatusNotFound:
		log.Println("in not found handler")
		writer.WriteHeader(statusCode)
		writer.Header().Set("Content-Type", "application/json")
		m := map[string]string{
			"error": "not found",
		}
		b, _ := json.Marshal(m)
		_, _ = writer.Write(b)
	case http.StatusBadRequest:
		log.Println("in bad request handler")
	case http.StatusMethodNotAllowed:
		log.Println("method not allowed")
		writer.WriteHeader(statusCode)
	}
}

func errorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	var apiError *APIError
	ok := errors.As(err, &apiError)
	if ok {
		resp := gErrorResponse(err)
		b, _ := json.Marshal(resp)
		w.WriteHeader(apiError.StatusCode)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(b)
	}

	//var newError runtime.HTTPStatusError
	//s, ok := status.FromError(err)
	//if !ok {
	//	st, _ := strconv.Atoi(s.Code().String())
	//	newError.HTTPStatus = st
	//	newError.Err = err
	//}
	//runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, writer, request, &newError)
}
