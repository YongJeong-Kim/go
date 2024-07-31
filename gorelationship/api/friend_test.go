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
	mockfriendsvc "gorelationship/service/mock/friend"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFriendAccept(t *testing.T) {
	testCases := []struct {
		name       string
		param      map[string]string
		auth       string
		run        func(t *testing.T, auth, method, url string, p map[string]string) *http.Request
		buildStubs func(t *testing.T, m *mockfriendsvc.MockFriender, auth string, p map[string]string)
		check      func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			param: map[string]string{
				"from_request": "2529ceea-781a-48c1-9f5b-b689989974f4",
			},
			auth: uuid.NewString(),
			run: func(t *testing.T, auth, method, url string, p map[string]string) *http.Request {
				b, err := json.Marshal(p)
				require.NoError(t, err)
				req, err := http.NewRequest(method, url, bytes.NewReader(b))
				require.NoError(t, err)
				req.Header.Set("user", auth)
				return req
			},
			buildStubs: func(t *testing.T, m *mockfriendsvc.MockFriender, auth string, p map[string]string) {
				m.EXPECT().Accept(context.Background(), p["from_request"], auth).Times(1).Return(nil)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			name:  "empty from request or whitespace",
			param: map[string]string{},
			auth:  uuid.NewString(),
			run: func(t *testing.T, auth, method, url string, p map[string]string) *http.Request {
				b, err := json.Marshal(p)
				require.NoError(t, err)
				req, err := http.NewRequest(method, url, bytes.NewReader(b))
				require.NoError(t, err)
				req.Header.Set("user", auth)
				return req
			},
			buildStubs: func(t *testing.T, m *mockfriendsvc.MockFriender, auth string, p map[string]string) {
				m.EXPECT().Accept(context.Background(), gomock.Any(), auth).Times(0)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				body, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)
				var r map[string]string
				err = json.Unmarshal(body, &r)
				require.NoError(t, err)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Equal(t, "Key: 'FromRequest' Error:Field validation for 'FromRequest' failed on the 'required' tag", r["error"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s, mf, _ := newTestServer(t)
			tc.buildStubs(t, mf, tc.auth, tc.param)
			req := tc.run(t, tc.auth, http.MethodPut, "/friends/accept", tc.param)
			recorder := httptest.NewRecorder()
			s.Router.ServeHTTP(recorder, req)
			tc.check(t, recorder)
		})
	}
}

func TestFriendCount(t *testing.T) {
	testCases := []struct {
		name       string
		auth       string
		run        func(t *testing.T, auth, method, url string) *http.Request
		buildStubs func(t *testing.T, m *mockfriendsvc.MockFriender, auth string)
		check      func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			auth: uuid.NewString(),
			run: func(t *testing.T, auth, method, url string) *http.Request {
				req, err := http.NewRequest(method, url, nil)
				require.NoError(t, err)
				req.Header.Set("user", auth)
				return req
			},
			buildStubs: func(t *testing.T, m *mockfriendsvc.MockFriender, auth string) {
				m.EXPECT().Count(context.Background(), auth).Times(1).Return(int64(1), nil)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				body, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)
				var r map[string]int
				err = json.Unmarshal(body, &r)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Equal(t, 1, r["count"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s, mf, _ := newTestServer(t)
			tc.buildStubs(t, mf, tc.auth)
			req := tc.run(t, tc.auth, http.MethodGet, "/friends/count")
			recorder := httptest.NewRecorder()
			s.Router.ServeHTTP(recorder, req)
			tc.check(t, recorder)
		})
	}
}

func TestListFriends(t *testing.T) {
	testCases := []struct {
		name       string
		auth       string
		response   []repository.ListResult
		run        func(t *testing.T, auth, method, url string) *http.Request
		buildStubs func(t *testing.T, m *mockfriendsvc.MockFriender, auth string, resp []repository.ListResult)
		check      func(t *testing.T, recorder *httptest.ResponseRecorder, resp []repository.ListResult)
	}{
		{
			name: "OK",
			auth: uuid.NewString(),
			response: []repository.ListResult{
				{
					ID:          uuid.NewString(),
					Name:        "",
					CreatedDate: time.Now().UTC(),
				},
			},
			run: func(t *testing.T, auth, method, url string) *http.Request {
				req, err := http.NewRequest(method, url, nil)
				require.NoError(t, err)
				req.Header.Set("user", auth)
				return req
			},
			buildStubs: func(t *testing.T, m *mockfriendsvc.MockFriender, auth string, resp []repository.ListResult) {
				m.EXPECT().List(context.Background(), auth).Times(1).Return(resp, nil)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder, resp []repository.ListResult) {
				body, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)
				var r []repository.ListResult
				err = json.Unmarshal(body, &r)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Equal(t, 1, len(r))
				for i := range r {
					require.Equal(t, resp[i].ID, r[i].ID)
					require.Equal(t, resp[i].Name, r[i].Name)
					require.Equal(t, resp[i].CreatedDate, r[i].CreatedDate)
				}
			},
		},
		{
			name: "invalid user id",
			auth: "invalid user id",
			response: []repository.ListResult{
				{
					ID:          uuid.NewString(),
					Name:        "",
					CreatedDate: time.Now().UTC(),
				},
			},
			run: func(t *testing.T, auth, method, url string) *http.Request {
				req, err := http.NewRequest(method, url, nil)
				require.NoError(t, err)
				req.Header.Set("user", auth)
				return req
			},
			buildStubs: func(t *testing.T, m *mockfriendsvc.MockFriender, auth string, resp []repository.ListResult) {
				m.EXPECT().List(context.Background(), auth).Times(1).Return(nil, errors.New("invalid user uuid"))
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder, resp []repository.ListResult) {
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
			s, mf, _ := newTestServer(t)
			tc.buildStubs(t, mf, tc.auth, tc.response)
			req := tc.run(t, tc.auth, http.MethodGet, "/friends")
			recorder := httptest.NewRecorder()
			s.Router.ServeHTTP(recorder, req)
			tc.check(t, recorder, tc.response)
		})
	}
}

func TestMutualFriends(t *testing.T) {
	testCases := []struct {
		name       string
		auth       string
		param      map[string]string
		response   []repository.ListMutualsResult
		run        func(t *testing.T, auth, method, url string, p map[string]string) *http.Request
		buildStubs func(t *testing.T, m *mockfriendsvc.MockFriender, auth string, p map[string]string, resp []repository.ListMutualsResult)
		check      func(t *testing.T, recorder *httptest.ResponseRecorder, resp []repository.ListMutualsResult)
	}{
		{
			name: "OK",
			auth: uuid.NewString(),
			param: map[string]string{
				"user_id": uuid.NewString(),
			},
			response: []repository.ListMutualsResult{
				{
					ID:          uuid.NewString(),
					Name:        "qweqwe",
					CreatedDate: time.Now().UTC(),
				},
			},
			run: func(t *testing.T, auth, method, url string, p map[string]string) *http.Request {
				req, err := http.NewRequest(method, url+"/"+p["user_id"], nil)
				require.NoError(t, err)
				req.Header.Set("user", auth)
				return req
			},
			buildStubs: func(t *testing.T, m *mockfriendsvc.MockFriender, auth string, p map[string]string, resp []repository.ListMutualsResult) {
				m.EXPECT().ListMutuals(context.Background(), auth, p["user_id"]).Times(1).Return(resp, nil)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder, resp []repository.ListMutualsResult) {
				body, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)
				var r []repository.ListMutualsResult
				err = json.Unmarshal(body, &r)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, recorder.Code)
				for i := range r {
					require.Equal(t, resp[i].ID, r[i].ID)
					require.Equal(t, resp[i].Name, r[i].Name)
					require.Equal(t, resp[i].CreatedDate, r[i].CreatedDate)
				}
			},
		},
		{
			name: "invalid user id",
			auth: uuid.NewString(),
			param: map[string]string{
				"user_id": "invalid user id",
			},
			response: []repository.ListMutualsResult{},
			run: func(t *testing.T, auth, method, url string, p map[string]string) *http.Request {
				req, err := http.NewRequest(method, url+"/"+p["user_id"], nil)
				require.NoError(t, err)
				req.Header.Set("user", auth)
				return req
			},
			buildStubs: func(t *testing.T, m *mockfriendsvc.MockFriender, auth string, p map[string]string, resp []repository.ListMutualsResult) {
				m.EXPECT().ListMutuals(context.Background(), auth, gomock.Any()).Times(0)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder, resp []repository.ListMutualsResult) {
				body, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)
				var r map[string]string
				err = json.Unmarshal(body, &r)
				require.NoError(t, err)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Equal(t, "invalid UUID length: 15", r["error"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s, mf, _ := newTestServer(t)
			tc.buildStubs(t, mf, tc.auth, tc.param, tc.response)
			req := tc.run(t, tc.auth, http.MethodGet, "/friends/mutual", tc.param)
			recorder := httptest.NewRecorder()
			s.Router.ServeHTTP(recorder, req)
			tc.check(t, recorder, tc.response)
		})
	}
}

func TestListFromRequests(t *testing.T) {
	testCases := []struct {
		name       string
		auth       string
		response   []repository.ListFromRequestsResult
		run        func(t *testing.T, auth, method, url string) *http.Request
		buildStubs func(t *testing.T, m *mockfriendsvc.MockFriender, auth string, resp []repository.ListFromRequestsResult)
		check      func(t *testing.T, recorder *httptest.ResponseRecorder, resp []repository.ListFromRequestsResult)
	}{
		{
			name: "OK",
			auth: uuid.NewString(),
			response: []repository.ListFromRequestsResult{
				{
					ID:          uuid.NewString(),
					Name:        "qweqwe",
					CreatedDate: time.Now().UTC(),
				},
			},
			run: func(t *testing.T, auth, method, url string) *http.Request {
				req, err := http.NewRequest(method, url, nil)
				require.NoError(t, err)
				req.Header.Set("user", auth)
				return req
			},
			buildStubs: func(t *testing.T, m *mockfriendsvc.MockFriender, auth string, resp []repository.ListFromRequestsResult) {
				m.EXPECT().ListFromRequests(context.Background(), auth).Times(1).Return(resp, nil)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder, resp []repository.ListFromRequestsResult) {
				body, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)
				var r []repository.ListMutualsResult
				err = json.Unmarshal(body, &r)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, recorder.Code)
				for i := range r {
					require.Equal(t, resp[i].ID, r[i].ID)
					require.Equal(t, resp[i].Name, r[i].Name)
					require.Equal(t, resp[i].CreatedDate, r[i].CreatedDate)
				}
			},
		},
		{
			name:     "invalid user id",
			auth:     "invalid user id",
			response: []repository.ListFromRequestsResult{},
			run: func(t *testing.T, auth, method, url string) *http.Request {
				req, err := http.NewRequest(method, url, nil)
				require.NoError(t, err)
				req.Header.Set("user", auth)
				return req
			},
			buildStubs: func(t *testing.T, m *mockfriendsvc.MockFriender, auth string, resp []repository.ListFromRequestsResult) {
				m.EXPECT().ListFromRequests(context.Background(), auth).Times(1).Return(nil, errors.New("invalid user uuid"))
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder, resp []repository.ListFromRequestsResult) {
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
			s, mf, _ := newTestServer(t)
			tc.buildStubs(t, mf, tc.auth, tc.response)
			req := tc.run(t, tc.auth, http.MethodGet, "/friends/request")
			recorder := httptest.NewRecorder()
			s.Router.ServeHTTP(recorder, req)
			tc.check(t, recorder, tc.response)
		})
	}
}

func TestMutualFriendCount(t *testing.T) {
	testCases := []struct {
		name       string
		auth       string
		param      map[string]string
		response   int64
		run        func(t *testing.T, auth, method, url string, p map[string]string) *http.Request
		buildStubs func(t *testing.T, m *mockfriendsvc.MockFriender, auth string, p map[string]string, resp int64)
		check      func(t *testing.T, recorder *httptest.ResponseRecorder, resp int64)
	}{
		{
			name: "OK",
			auth: uuid.NewString(),
			param: map[string]string{
				"user_id": uuid.NewString(),
			},
			response: 1,
			run: func(t *testing.T, auth, method, url string, p map[string]string) *http.Request {
				req, err := http.NewRequest(method, url+"/"+p["user_id"]+"/count", nil)
				require.NoError(t, err)
				req.Header.Set("user", auth)
				return req
			},
			buildStubs: func(t *testing.T, m *mockfriendsvc.MockFriender, auth string, p map[string]string, resp int64) {
				m.EXPECT().MutualCount(context.Background(), auth, p["user_id"]).Times(1).Return(resp, nil)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder, resp int64) {
				body, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)
				var r map[string]int64
				err = json.Unmarshal(body, &r)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Equal(t, resp, r["count"])
			},
		},
		{
			name: "invalid yourself user id",
			auth: "invalid user id",
			param: map[string]string{
				"user_id": uuid.NewString(),
			},
			response: 1,
			run: func(t *testing.T, auth, method, url string, p map[string]string) *http.Request {
				req, err := http.NewRequest(method, url+"/"+p["user_id"]+"/count", nil)
				require.NoError(t, err)
				req.Header.Set("user", auth)
				return req
			},
			buildStubs: func(t *testing.T, m *mockfriendsvc.MockFriender, auth string, p map[string]string, resp int64) {
				m.EXPECT().MutualCount(context.Background(), auth, gomock.Any()).Times(1).Return(int64(0), errors.New("invalid user1 uuid"))
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder, resp int64) {
				body, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)
				var r map[string]string
				err = json.Unmarshal(body, &r)
				require.NoError(t, err)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Equal(t, "invalid user1 uuid", r["error"])
			},
		},
		{
			name: "invalid user id",
			auth: uuid.NewString(),
			param: map[string]string{
				"user_id": "invalid user id",
			},
			response: 1,
			run: func(t *testing.T, auth, method, url string, p map[string]string) *http.Request {
				req, err := http.NewRequest(method, url+"/"+p["user_id"]+"/count", nil)
				require.NoError(t, err)
				req.Header.Set("user", auth)
				return req
			},
			buildStubs: func(t *testing.T, m *mockfriendsvc.MockFriender, auth string, p map[string]string, resp int64) {
				m.EXPECT().MutualCount(context.Background(), auth, gomock.Any()).Times(1).Return(int64(0), errors.New("invalid user2 uuid"))
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder, resp int64) {
				body, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)
				var r map[string]string
				err = json.Unmarshal(body, &r)
				require.NoError(t, err)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Equal(t, "invalid user2 uuid", r["error"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s, mf, _ := newTestServer(t)
			tc.buildStubs(t, mf, tc.auth, tc.param, tc.response)
			req := tc.run(t, tc.auth, http.MethodGet, "/friends/mutual", tc.param)
			recorder := httptest.NewRecorder()
			s.Router.ServeHTTP(recorder, req)
			tc.check(t, recorder, tc.response)
		})
	}
}

func TestFriendRequest(t *testing.T) {
	testCases := []struct {
		name       string
		param      map[string]string
		auth       string
		run        func(t *testing.T, auth, method, url string, p map[string]string) *http.Request
		buildStubs func(t *testing.T, m *mockfriendsvc.MockFriender, auth string, p map[string]string)
		check      func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			param: map[string]string{
				"user_id": "2529ceea-781a-48c1-9f5b-b689989974f4",
			},
			auth: uuid.NewString(),
			run: func(t *testing.T, auth, method, url string, p map[string]string) *http.Request {
				b, err := json.Marshal(p)
				require.NoError(t, err)
				req, err := http.NewRequest(method, url, bytes.NewReader(b))
				require.NoError(t, err)
				req.Header.Set("user", auth)
				return req
			},
			buildStubs: func(t *testing.T, m *mockfriendsvc.MockFriender, auth string, p map[string]string) {
				m.EXPECT().Request(context.Background(), auth, p["user_id"]).Times(1).Return(nil)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name:  "empty to request user id",
			param: map[string]string{},
			auth:  uuid.NewString(),
			run: func(t *testing.T, auth, method, url string, p map[string]string) *http.Request {
				b, err := json.Marshal(p)
				require.NoError(t, err)
				req, err := http.NewRequest(method, url, bytes.NewReader(b))
				require.NoError(t, err)
				req.Header.Set("user", auth)
				return req
			},
			buildStubs: func(t *testing.T, m *mockfriendsvc.MockFriender, auth string, p map[string]string) {
				m.EXPECT().Request(context.Background(), auth, gomock.Any()).Times(0)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				body, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)
				var r map[string]string
				err = json.Unmarshal(body, &r)
				require.NoError(t, err)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Equal(t, "Key: 'UserID' Error:Field validation for 'UserID' failed on the 'required' tag", r["error"])
			},
		},
		{
			name: "request yourself",
			param: map[string]string{
				"user_id": "78ded853-1017-45ca-bfc8-5e5103ccfa88",
			},
			auth: "78ded853-1017-45ca-bfc8-5e5103ccfa88",
			run: func(t *testing.T, auth, method, url string, p map[string]string) *http.Request {
				b, err := json.Marshal(p)
				require.NoError(t, err)
				req, err := http.NewRequest(method, url, bytes.NewReader(b))
				require.NoError(t, err)
				req.Header.Set("user", auth)
				return req
			},
			buildStubs: func(t *testing.T, m *mockfriendsvc.MockFriender, auth string, p map[string]string) {
				m.EXPECT().Request(context.Background(), auth, p["user_id"]).Times(1).Return(errors.New("cannot request yourself"))
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				body, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)
				var r map[string]string
				err = json.Unmarshal(body, &r)
				require.NoError(t, err)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Equal(t, "cannot request yourself", r["error"])
			},
		},
		{
			name: "invalid yourself user id",
			param: map[string]string{
				"user_id": uuid.NewString(),
			},
			auth: "invalid user id",
			run: func(t *testing.T, auth, method, url string, p map[string]string) *http.Request {
				b, err := json.Marshal(p)
				require.NoError(t, err)
				req, err := http.NewRequest(method, url, bytes.NewReader(b))
				require.NoError(t, err)
				req.Header.Set("user", auth)
				return req
			},
			buildStubs: func(t *testing.T, m *mockfriendsvc.MockFriender, auth string, p map[string]string) {
				m.EXPECT().Request(context.Background(), auth, p["user_id"]).Times(1).Return(errors.New("invalid request user uuid"))
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				body, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)
				var r map[string]string
				err = json.Unmarshal(body, &r)
				require.NoError(t, err)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Equal(t, "invalid request user uuid", r["error"])
			},
		},
		{
			name: "invalid user id",
			param: map[string]string{
				"user_id": uuid.NewString(),
			},
			auth: uuid.NewString(),
			run: func(t *testing.T, auth, method, url string, p map[string]string) *http.Request {
				b, err := json.Marshal(p)
				require.NoError(t, err)
				req, err := http.NewRequest(method, url, bytes.NewReader(b))
				require.NoError(t, err)
				req.Header.Set("user", auth)
				return req
			},
			buildStubs: func(t *testing.T, m *mockfriendsvc.MockFriender, auth string, p map[string]string) {
				m.EXPECT().Request(context.Background(), auth, p["user_id"]).Times(1).Return(errors.New("invalid accept user uuid"))
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				body, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)
				var r map[string]string
				err = json.Unmarshal(body, &r)
				require.NoError(t, err)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Equal(t, "invalid accept user uuid", r["error"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s, mf, _ := newTestServer(t)
			tc.buildStubs(t, mf, tc.auth, tc.param)
			req := tc.run(t, tc.auth, http.MethodPost, "/friends/request", tc.param)
			recorder := httptest.NewRecorder()
			s.Router.ServeHTTP(recorder, req)
			tc.check(t, recorder)
		})
	}
}

func TestListFromRequestCount(t *testing.T) {
	testCases := []struct {
		name       string
		auth       string
		run        func(t *testing.T, auth, method, url string) *http.Request
		buildStubs func(t *testing.T, m *mockfriendsvc.MockFriender, auth string)
		check      func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			auth: uuid.NewString(),
			run: func(t *testing.T, auth, method, url string) *http.Request {
				req, err := http.NewRequest(method, url, nil)
				require.NoError(t, err)
				req.Header.Set("user", auth)
				return req
			},
			buildStubs: func(t *testing.T, m *mockfriendsvc.MockFriender, auth string) {
				m.EXPECT().FromRequestCount(context.Background(), auth).Times(1).Return(int64(1), nil)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				body, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)
				var r map[string]int64
				err = json.Unmarshal(body, &r)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Equal(t, int64(1), r["count"])
			},
		},
		{
			name: "invalid user id",
			auth: "invalid user id",
			run: func(t *testing.T, auth, method, url string) *http.Request {
				req, err := http.NewRequest(method, url, nil)
				require.NoError(t, err)
				req.Header.Set("user", auth)
				return req
			},
			buildStubs: func(t *testing.T, m *mockfriendsvc.MockFriender, auth string) {
				m.EXPECT().FromRequestCount(context.Background(), auth).Times(1).Return(int64(0), errors.New("invalid user uuid"))
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
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
			s, mf, _ := newTestServer(t)
			tc.buildStubs(t, mf, tc.auth)
			req := tc.run(t, tc.auth, http.MethodGet, "/friends/request/count")
			recorder := httptest.NewRecorder()
			s.Router.ServeHTTP(recorder, req)
			tc.check(t, recorder)
		})
	}
}
