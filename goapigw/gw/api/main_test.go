package api

import "github.com/yongjeong-kim/go/goapigw/gw/token"

func newTestServer(tokenVerifier token.TokenVerifier) *Server {
	return &Server{
		TokenVerifier: tokenVerifier,
	}
}
