package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/yongjeong-kim/go/goapigw/gw/token"
	tkmock "github.com/yongjeong-kim/go/goapigw/gw/token/mock"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name       string
		addAuth    func(*http.Request)
		buildStubs func(*tkmock.MockTokenVerifier)
		check      func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "empty auth",
			addAuth: func(req *http.Request) {
			},
			buildStubs: func(m *tkmock.MockTokenVerifier) {
				m.EXPECT().Verify(gomock.Any()).Times(0)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "OK",
			addAuth: func(req *http.Request) {
				auth := fmt.Sprintf("%s %s", authorizationTypeBearer, "v4.local.1rgk1mWhxlIZTibar2LtKRHAe_Vz-dalZVpa0SQ4EtmPh9XOVOmm5W_QckylWvn-2m9q294_Y0DQ8vu9r8Ete7u9BAKis6Gb2Vvy7GoVWRWSJ4y2kmPYtGRpzmBTOURtP6OaFwVAb6yQDjk4gR7-KNVmxSO32NSck0PTwYkvPYs9zv6ZZhBST31ij7GXrfOUzqI8cbV9ftaCmj_-IjHbLbIj4C-y9v480C8mbl5sLgS2X9in5WIivA")
				req.Header.Set(authorizationHeaderKey, auth)
			},
			buildStubs: func(m *tkmock.MockTokenVerifier) {
				m.EXPECT().Verify(gomock.Any()).Times(1).Return(&token.Payload{}, nil)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "no bearer",
			addAuth: func(req *http.Request) {
				auth := fmt.Sprintf("%s %s", "no bearer", "v4.local.1rgk1mWhxlIZTibar2LtKRHAe_Vz-dalZVpa0SQ4EtmPh9XOVOmm5W_QckylWvn-2m9q294_Y0DQ8vu9r8Ete7u9BAKis6Gb2Vvy7GoVWRWSJ4y2kmPYtGRpzmBTOURtP6OaFwVAb6yQDjk4gR7-KNVmxSO32NSck0PTwYkvPYs9zv6ZZhBST31ij7GXrfOUzqI8cbV9ftaCmj_-IjHbLbIj4C-y9v480C8mbl5sLgS2X9in5WIivA")
				req.Header.Set(authorizationHeaderKey, auth)
			},
			buildStubs: func(m *tkmock.MockTokenVerifier) {
				m.EXPECT().Verify(gomock.Any()).Times(0)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				b, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				ce := make(map[string]any)
				err = json.Unmarshal(b, &ce)
				require.NoError(t, err)
				require.Equal(t, "unsupported authorization type "+"no", (ce["inner"]).(string))
				require.Equal(t, http.StatusUnauthorized, int((ce["status_code"]).(float64)))
				require.Equal(t, "unsupported authorization type "+"no", (ce["message"]).(string))
			},
		},
		{
			name: "empty bearer",
			addAuth: func(req *http.Request) {
				auth := fmt.Sprintf("%s %s", "", "v4.local.1rgk1mWhxlIZTibar2LtKRHAe_Vz-dalZVpa0SQ4EtmPh9XOVOmm5W_QckylWvn-2m9q294_Y0DQ8vu9r8Ete7u9BAKis6Gb2Vvy7GoVWRWSJ4y2kmPYtGRpzmBTOURtP6OaFwVAb6yQDjk4gR7-KNVmxSO32NSck0PTwYkvPYs9zv6ZZhBST31ij7GXrfOUzqI8cbV9ftaCmj_-IjHbLbIj4C-y9v480C8mbl5sLgS2X9in5WIivA")
				req.Header.Set(authorizationHeaderKey, auth)
			},
			buildStubs: func(m *tkmock.MockTokenVerifier) {
				m.EXPECT().Verify(gomock.Any()).Times(0)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				b, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				ce := make(map[string]any)
				err = json.Unmarshal(b, &ce)
				require.NoError(t, err)
				require.Equal(t, "invalid authorization header format", (ce["inner"]).(string))
				require.Equal(t, http.StatusUnauthorized, int((ce["status_code"]).(float64)))
				require.Equal(t, "invalid authorization header format", (ce["message"]).(string))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := tkmock.NewMockTokenVerifier(ctrl)
			server := newTestServer(m)
			setupRouterForTest(t, server)
			server.Router.Use(authMiddleware(server.TokenVerifier)).GET("/auth", func(c *gin.Context) {
				c.Status(http.StatusOK)
				return
			})
			tc.buildStubs(m)
			req, err := http.NewRequest(http.MethodGet, "/auth", nil)
			require.NoError(t, err)
			tc.addAuth(req)
			recorder := httptest.NewRecorder()

			server.Router.ServeHTTP(recorder, req)
			tc.check(t, recorder)
		})
	}
}
