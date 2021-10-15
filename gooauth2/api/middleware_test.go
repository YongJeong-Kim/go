package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gooauth2/token"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func addAuthorization(
	t *testing.T,
	request *http.Request,
	server *Server,
	authorizationType string,
	username string,
	duration token.JWTDuration,
) {
	payloadDetails, err := server.maker.CreateToken(username, duration)
	require.NoError(t, err)

	err = server.CreateAuth(payloadDetails.Payload.AccessTokenPayload.UserID, payloadDetails)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, payloadDetails.Token["access_token"])
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, server *Server)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				duration := token.JWTDuration{
					AccessTokenDuration:  time.Minute,
					RefreshTokenDuration: time.Minute,
				}
				addAuthorization(t, request, server, authorizationTypeBearer, "authorized username", duration)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Anonymous Unauthorized",
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Incorrect Access Token Unauthorized",
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				duration := token.JWTDuration{
					AccessTokenDuration:  time.Minute,
					RefreshTokenDuration: time.Minute,
				}
				addAuthorization(t, request, server, authorizationTypeBearer, "authorized username", duration)
				request.Header.Set(authorizationHeaderKey, "incorrent access token")
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Expired Access Token",
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				duration := token.JWTDuration{
					AccessTokenDuration:  -time.Minute,
					RefreshTokenDuration: -time.Minute,
				}
				addAuthorization(t, request, server, authorizationTypeBearer, "authorized username", duration)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t)
			authPath := "/authmiddlewaretest"
			server.router.GET(
				authPath,
				server.authMiddleware(server.maker),
				func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
