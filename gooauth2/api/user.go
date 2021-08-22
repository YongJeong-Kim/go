package api

import (
	"gooauth2/token"
	"time"
)

func CreateAuth(payload token.AccessTokenPayload) {
	at := payload.ExpiredAt
	rt := payload.RtExpiredAt
	now := time.Now()
}
