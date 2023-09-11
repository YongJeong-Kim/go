package api

import (
	mocksrv "github.com/YongJeong-Kim/go/goasynq/service/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mocksrv.NewMockServicer(ctrl)
	mock.EXPECT().CreateUser(gomock.Any()).Times(1).Return(nil)

	server := NewServer(mock, nil, nil)
	server.SetupRouter()

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/users", nil)
	require.NoError(t, err)

	server.Router.ServeHTTP(recorder, req)

	require.Equal(t, recorder.Code, http.StatusCreated)
}
