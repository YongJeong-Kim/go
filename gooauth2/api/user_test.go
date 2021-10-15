package api

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"gooauth2/token"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func addRefreshTokenOnly(
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

	refreshToken := fmt.Sprintf("%s %s", authorizationType, payloadDetails.Token["refresh_token"])
	request.Header.Set(authorizationHeaderKey, refreshToken)
}

func TestUser(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, server *Server)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "NewAccessTokenProvided",
			setupAuth: func(t *testing.T, request *http.Request, server *Server) {
				duration := token.JWTDuration{
					AccessTokenDuration:  -time.Minute,
					RefreshTokenDuration: time.Minute,
				}
				addRefreshTokenOnly(t, request, server, authorizationTypeBearer, "authorized username", duration)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t)

			recorder := httptest.NewRecorder()
			url := "/token/refresh"
			request, err := http.NewRequest(http.MethodPost, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
