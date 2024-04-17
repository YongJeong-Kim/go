package gapi

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/protobuf/proto"
	"io"
	"mime/multipart"
	"net/url"
)

var _ runtime.Marshaler = (*MultipartFormPb)(nil)

type MultipartFormPb struct {
	runtime.Marshaler
}

func GWMultipartForm(marshaler runtime.Marshaler) runtime.ServeMuxOption {
	return runtime.WithMarshalerOption("multipart/form-data", &MultipartFormPb{
		Marshaler: marshaler,
	})
}

func (j *MultipartFormPb) NewDecoder(r io.Reader) runtime.Decoder {
	return runtime.DecoderFunc(func(v any) error {
		msg, ok := v.(proto.Message)
		if !ok {
			return fmt.Errorf("not proto message") // nolint:goerr113
		}

		br := bufio.NewReaderSize(r, 1024)
		pb, err := br.Peek(100)
		if err != nil {
			return fmt.Errorf("peek boundary: %w", err)
		}

		if len(pb) < 2 {
			return fmt.Errorf("boundary len < 2") // nolint:goerr113
		}

		boundary := bytes.TrimSpace(bytes.Split(pb, []byte("\n"))[0])[2:]

		values := make(url.Values)

		mr := multipart.NewReader(br, string(boundary))
		for {
			var p *multipart.Part
			p, err = mr.NextPart()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				return fmt.Errorf("read next part: %w", err)
			}

			var data []byte
			data, err = io.ReadAll(p)
			if err != nil {
				return fmt.Errorf("read part body: %w", err)
			}

			if p.FileName() != "" {
				/*
					in proto file:
						message Media {
							string filename = 1;
							string content_type = 2;
							bytes content = 3;
						}
				*/
				values.Set(p.FormName()+".filename", p.FileName())
				values.Set(p.FormName()+".content_type", p.Header.Get("Content-Type"))
				values.Set(p.FormName()+".content", base64.StdEncoding.EncodeToString(data))
			} else {
				values.Set(p.FormName(), string(data))
			}
		}

		err = runtime.PopulateQueryParameters(msg, values, &utilities.DoubleArray{})
		if err != nil {
			return fmt.Errorf("papulate query params: %w", err)
		}

		return nil
	})
}

func (j *MultipartFormPb) Unmarshal(data []byte, v any) error {
	return j.NewDecoder(bytes.NewReader(data)).Decode(v) // nolint: wrapcheck
}
