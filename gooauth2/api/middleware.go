package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gooauth2/token"
	"net/http"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func (server *Server) authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		// if len(authorizationHeader) == 0 {
		// 	err := errors.New("authorization header is not provided")
		// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		// 	return
		// }

		// fields := strings.Fields(authorizationHeader)
		// if len(fields) < 2 {
		// 	err := errors.New("invalid authorization header format")
		// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		// 	return
		// }

		// authorizationType := strings.ToLower(fields[0])
		// if authorizationType != authorizationTypeBearer {
		// 	err := fmt.Errorf("unsupported authorization type %s", authorizationType)
		// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		// 	return
		// }

		accessToken, err := server.maker.ExtractToken(authorizationHeader)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		payload, err := tokenMaker.VerifyAccessToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessTokenID := payload.ID
		_, err = server.fetchAuth(accessTokenID.String())
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

func (server *Server) fetchAuth(accessToken string) (string, error) {
	ctx := context.Background()
	userID, err := server.redisClient.Get(ctx, accessToken).Result()
	if err != nil {
		if err == redis.Nil {
			return "", err
		}
		return "", err
	}

	return userID, err
}
