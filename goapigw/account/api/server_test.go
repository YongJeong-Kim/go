package api

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	svcmock "github.com/yongjeong-kim/go/goapigw/account/service/mock"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	url "net/url"
	"strings"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	testCases := []struct {
		name       string
		buildStubs func(*svcmock.MockAccountServicer)
		check      func(*httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStubs: func(m *svcmock.MockAccountServicer) {
				m.EXPECT().Login("aaa", "1234", time.Minute).Times(1).Return("asdf", nil)
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := svcmock.NewMockAccountServicer(ctrl)
			server := newTestServer(m)
			server.SetupRouter()
			tc.buildStubs(m)

			form := url.Values{}
			form.Add("username", "aaa")
			form.Add("password", "1234")

			loginURL := Accountv1 + "/login"
			request, err := http.NewRequest(http.MethodPost, loginURL, strings.NewReader(form.Encode()))
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			recorder := httptest.NewRecorder()
			server.Router.ServeHTTP(recorder, request)

			tc.check(recorder)
		})
	}
}
