package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	token "github.com/huzaifa678/Crypto-currency-web-app-project/token"
)


const (
	AuthoriationHeaderKey = "authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)


func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(AuthoriationHeaderKey)

		if len(authHeader) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"unauthorized": "authorization header is not provided"})
			return
		}

		fields := strings.Fields(authHeader)

		if len(fields) < 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"unauthorized": "invalid format for the authorization header"})
			return
		}

		authorizationTypeBearer := strings.ToLower(fields[0])
		if authorizationTypeBearer != AuthorizationTypeBearer {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"unauthorized": "unsupported authorization type"})
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"unauthorized": "access token is not valid"})
			return
		}

		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}