package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorelationship/service"
	mockusersvc "gorelationship/service/mock/user"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		name     string
		param    map[string]string
		response map[string]string
		mock     func(t *testing.T, m *mockusersvc.MockUserManager, p map[string]string, createdID string)
		check    func(t *testing.T, recorder *httptest.ResponseRecorder, createdID string)
	}{
		{
			name: "OK",
			param: map[string]string{
				"name": "test",
			},
			mock: func(t *testing.T, m *mockusersvc.MockUserManager, p map[string]string, createdID string) {
				//m.EXPECT().Create(context.Background(), gomock.Eq(p["name"])).Times(1).Return(createdID, nil)
				m.EXPECT().Create(context.Background(), gomock.Any()).Times(1).Return(createdID, nil)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder, createdID string) {
				require.Equal(t, createdID, recorder.Body.String())
				require.Equal(t, http.StatusCreated, recorder.Code)
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
			tc.mock(t, mock, tc.param, createdID)

			svc := service.NewService(nil, mock)
			s := NewServer(svc)
			s.SetupRouter()

			b, err := json.Marshal(tc.param)
			require.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(b))
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			s.Router.ServeHTTP(recorder, req)
			tc.check(t, recorder, createdID)
		})
	}
}
