package token

type Maker interface {
	CreateToken(username string, tokenDuration JWTDuration) (*PayloadDetails, error)
	VerifyToken(payload string) (*Payload, error)
}
