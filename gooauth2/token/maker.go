package token

type Maker interface {
	CreateToken(username string, tokenDuration JWTDuration) (*PayloadDetails, error)
	VerifyToken(accessToken string) (*AccessTokenPayload, error)
	ExtractToken(authorizationHeader string) (string, error)
}
