package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorelationship/repository"
	"gorelationship/service"
	mockusersvc "gorelationship/service/mock/user"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		name       string
		param      map[string]string
		run        func(t *testing.T, p map[string]string, method, url string) *http.Request
		buildStubs func(t *testing.T, m *mockusersvc.MockUserManager, p map[string]string, createdID string)
		check      func(t *testing.T, recorder *httptest.ResponseRecorder, createdID string)
	}{
		{
			name: "OK",
			param: map[string]string{
				"name": "test",
			},
			run: func(t *testing.T, p map[string]string, method, url string) *http.Request {
				b, err := json.Marshal(p)
				require.NoError(t, err)
				req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(b))
				require.NoError(t, err)
				return req
			},
			buildStubs: func(t *testing.T, m *mockusersvc.MockUserManager, p map[string]string, createdID string) {
				m.EXPECT().Create(context.Background(), gomock.Any()).Times(1).Return(createdID, nil)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder, createdID string) {
				require.Equal(t, createdID, recorder.Body.String())
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name:  "no name",
			param: nil,
			run: func(t *testing.T, p map[string]string, method, url string) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users", nil)
				require.NoError(t, err)
				return req
			},
			buildStubs: func(t *testing.T, m *mockusersvc.MockUserManager, p map[string]string, createdID string) {
				m.EXPECT().Create(context.Background(), gomock.Any()).Times(0)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder, createdID string) {
				b := recorder.Body.Bytes()
				var r map[string]string
				err := json.Unmarshal(b, &r)
				require.NoError(t, err)

				require.Equal(t, "invalid request", r["error"])
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := mockusersvc.NewMockUserManager(ctrl)
			createdID := uuid.NewString()
			tc.buildStubs(t, mock, tc.param, createdID)

			svc := service.NewService(nil, mock)
			s := NewServer(svc)
			s.SetupRouter()

			req := tc.run(t, tc.param, http.MethodGet, "/users")
			recorder := httptest.NewRecorder()
			s.Router.ServeHTTP(recorder, req)
			tc.check(t, recorder, createdID)
		})
	}
}

func TestGetUser(t *testing.T) {
	testCases := []struct {
		name       string
		param      map[string]string
		result     *repository.GetResult
		run        func(t *testing.T, p map[string]string, method, url string) *http.Request
		buildStubs func(t *testing.T, m *mockusersvc.MockUserManager, p map[string]string, result *repository.GetResult)
		check      func(t *testing.T, recorder *httptest.ResponseRecorder, result *repository.GetResult)
	}{
		{
			name: "OK",
			param: map[string]string{
				"user_id": "c1c2fd15-fa55-4489-ad5c-78d1d779930f",
			},
			result: &repository.GetResult{
				ID:          "c1c2fd15-fa55-4489-ad5c-78d1d779930f",
				Name:        "asxz",
				CreatedDate: time.Now().UTC(),
			},
			run: func(t *testing.T, p map[string]string, method, url string) *http.Request {
				req, err := http.NewRequest(method, url+"/"+p["user_id"], nil)
				require.NoError(t, err)
				return req
			},
			buildStubs: func(t *testing.T, m *mockusersvc.MockUserManager, p map[string]string, result *repository.GetResult) {
				m.EXPECT().Get(context.Background(), p["user_id"]).Times(1).Return(result, nil)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder, result *repository.GetResult) {
				body, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)
				var r repository.GetResult
				err = json.Unmarshal(body, &r)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Equal(t, r.ID, result.ID)
				require.Equal(t, r.Name, result.Name)
				require.Equal(t, r.CreatedDate, result.CreatedDate)
			},
		},
		{
			name: "invalid user id",
			param: map[string]string{
				"user_id": "invalid user id",
			},
			result: &repository.GetResult{
				ID:          "any",
				Name:        "asxz",
				CreatedDate: time.Now().UTC(),
			},
			run: func(t *testing.T, p map[string]string, method, url string) *http.Request {
				req, err := http.NewRequest(method, url+"/"+p["user_id"], nil)
				require.NoError(t, err)
				return req
			},
			buildStubs: func(t *testing.T, m *mockusersvc.MockUserManager, p map[string]string, result *repository.GetResult) {
				m.EXPECT().Get(context.Background(), p["user_id"]).Times(1).Return(nil, errors.New("invalid user uuid"))
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder, result *repository.GetResult) {
				body, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)
				var r map[string]string
				err = json.Unmarshal(body, &r)
				require.NoError(t, err)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Equal(t, "invalid user uuid", r["error"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			svr, _, mu := newTestServer(t)
			tc.buildStubs(t, mu, tc.param, tc.result)
			req := tc.run(t, tc.param, http.MethodGet, "/users")
			recorder := httptest.NewRecorder()
			svr.Router.ServeHTTP(recorder, req)
			tc.check(t, recorder, tc.result)
		})
	}
}
