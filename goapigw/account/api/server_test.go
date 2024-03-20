package api

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	svcmock "github.com/yongjeong-kim/go/goapigw/account/service/mock"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	testCases := []struct {
		name       string
		param      url.Values
		buildStubs func(*svcmock.MockAccountServicer, url.Values)
		check      func(*httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			param: url.Values{
				"username": {"aaa"},
				"password": {"1234"},
			},
			buildStubs: func(m *svcmock.MockAccountServicer, u url.Values) {
				m.EXPECT().Login(u.Get("username"), u.Get("password"), 30*24*time.Hour).Times(1).Return("asdf", nil)
			},
			check: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				b, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				result := make(map[string]string)
				err = json.Unmarshal(b, &result)
				require.NoError(t, err)
				require.Equal(t, result["access_token"], "asdf")
			},
		},
		{
			name: "invalid param",
			param: url.Values{
				"username1212": {"aaa"},
				"password":     {"1234"},
			},
			buildStubs: func(m *svcmock.MockAccountServicer, u url.Values) {
				m.EXPECT().Login(gomock.Any(), gomock.Any(), 30*24*time.Hour).Times(0)
			},
			check: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				b, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				resp := make(map[string]any)
				err = json.Unmarshal(b, &resp)
				require.NoError(t, err)

				require.Equal(t, "invalid username or password", resp["message"].(string))
				require.Equal(t, http.StatusBadRequest, int(resp["status_code"].(float64)))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := svcmock.NewMockAccountServicer(ctrl)
			server := newTestServer(m)
			server.SetupRouter()
			tc.buildStubs(m, tc.param)

			loginURL := Accountv1 + "/login"
			request, err := http.NewRequest(http.MethodPost, loginURL, strings.NewReader(tc.param.Encode()))
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			recorder := httptest.NewRecorder()
			server.Router.ServeHTTP(recorder, request)

			tc.check(recorder)
		})
	}
}
