package token

type Maker interface {
	CreateToken(username string, tokenDuration JWTDuration) (map[string]string, error)
	VerifyToken(payload string) (*Payload, error)
}
