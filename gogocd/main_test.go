package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	code := m.Run()
	os.Exit(code)
}

func TestDefault(t *testing.T) {
	testCases := []struct {
		name          string
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.True(t, false)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			s := NewServer(gin.New())
			s.SetupRouter()

			recorder := httptest.NewRecorder()
			s.Router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}
