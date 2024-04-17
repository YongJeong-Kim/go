package gapi

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	svcmock "github.com/yongjeong-kim/go/goapigw/account/service/mock"
	accountv1 "github.com/yongjeong-kim/go/goapigw/accountpb/pb/account/v1"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLoginHTTP(t *testing.T) {
	testCases := []struct {
		name       string
		param      *accountv1.LoginRequest
		buildStubs func(*svcmock.MockAccountServicer, *accountv1.LoginRequest)
		check      func(*httptest.ResponseRecorder, *http.Request)
	}{
		{
			name: "OK",
			param: &accountv1.LoginRequest{
				Username: "aaa",
				Password: "1234",
			},
			buildStubs: func(m *svcmock.MockAccountServicer, param *accountv1.LoginRequest) {
				m.EXPECT().Login(param.GetUsername(), param.GetPassword(), 5*time.Minute).Times(1).Return("aa", nil)
			},
			check: func(recorder *httptest.ResponseRecorder, req *http.Request) {
				require.Equal(t, http.StatusOK, recorder.Code)
				b, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var resp map[string]string
				err = json.Unmarshal(b, &resp)
				require.NoError(t, err)

				require.Equal(t, "aa", resp["token"])
			},
		},
		{
			name: "invalid username or password",
			param: &accountv1.LoginRequest{
				Username: "asd",
				Password: "fe",
			},
			buildStubs: func(m *svcmock.MockAccountServicer, param *accountv1.LoginRequest) {
				m.EXPECT().Login(param.GetUsername(), param.GetPassword(), 5*time.Minute).Times(0)
			},
			check: func(recorder *httptest.ResponseRecorder, req *http.Request) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := svcmock.NewMockAccountServicer(ctrl)
			s := newTestServer(m)

			tc.buildStubs(m, tc.param)

			b, err := json.Marshal(tc.param)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/v1/account/login", bytes.NewReader(b))
			require.NoError(t, err)

			grpcMux := s.setupGatewayServer(context.Background())

			recorder := httptest.NewRecorder()
			s.GetRouter(grpcMux).ServeHTTP(recorder, req)

			tc.check(recorder, req)
		})
	}
}

func TestLoginGRPC(t *testing.T) {
	testCases := []struct {
		name       string
		param      *accountv1.LoginRequest
		buildStubs func(*svcmock.MockAccountServicer, *accountv1.LoginRequest)
		check      func(*httptest.ResponseRecorder, *http.Request)
	}{
		{
			name: "OK",
			param: &accountv1.LoginRequest{
				Username: "aaa",
				Password: "1234",
			},
			buildStubs: func(m *svcmock.MockAccountServicer, param *accountv1.LoginRequest) {
				m.EXPECT().Login(param.GetUsername(), param.GetPassword(), 5*time.Minute).Times(1).Return("aa", nil)
			},
			check: func(recorder *httptest.ResponseRecorder, req *http.Request) {

			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := svcmock.NewMockAccountServicer(ctrl)
			s := newTestServer(m)

			grpcServer := s.setupGRPCServer(context.Background())
			b, err := json.Marshal(tc.param)
			require.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/v1/account/login", bytes.NewReader(b))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/grpc")
			recorder := httptest.NewRecorder()
			grpcServer.ServeHTTP(recorder, req)
			tc.check(recorder, req)
		})
	}
}
