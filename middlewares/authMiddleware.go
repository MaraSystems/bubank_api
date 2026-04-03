package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/MaraSystems/bubank_api/utils"
	"github.com/gin-gonic/gin"
)

var (
	AuthHeaderKey      = "Authorization"
	AuthKey            = "auth"
	AuthType           = "Bearer"
	ErrAuthRequired    = fmt.Errorf("Authorization is required")
	ErrInvalidAuth     = fmt.Errorf("Authorization is invalid")
	ErrInvalidAuthType = fmt.Errorf("Authorization type is invalid")
)

func AuthMiddleWare(tokenMaker *utils.TokenMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.GetHeader(AuthHeaderKey)
		if authorization == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(ErrAuthRequired))
			return
		}

		fields := strings.Fields(authorization)
		if len(fields) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(ErrInvalidAuth))
			return
		}

		if fields[0] != AuthType {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(ErrInvalidAuthType))
			return
		}

		token := fields[1]
		payload, err := tokenMaker.Validate(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}

		ctx.Set(AuthKey, payload)
		ctx.Next()
	}
}
