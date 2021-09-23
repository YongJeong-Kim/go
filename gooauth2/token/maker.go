package token

type Maker interface {
	CreateToken(username string, tokenDuration JWTDuration) (*PayloadDetails, error)
	VerifyAccessToken(accessToken string) (*AccessTokenPayload, error)
	VerifyRefreshToken(refreshToken string) (*RefreshTokenPayload, error)
	ExtractToken(authorizationHeader string) (string, error)
}
