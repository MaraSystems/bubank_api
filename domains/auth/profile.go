package auth

import (
	"fmt"
	"net/http"

	"github.com/MaraSystems/bubank_api/domains/users"
	"github.com/MaraSystems/bubank_api/middlewares"
	"github.com/MaraSystems/bubank_api/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func (h AuthHandler) getProfile(ctx *gin.Context) {
	authorization := ctx.MustGet(middlewares.AuthKey).(*utils.TokenPayload)

	user, err := h.server.Store.GetUser(ctx.Request.Context(), authorization.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(fmt.Errorf("profile not found")))
			return
		}

		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	rsp := users.UserToHTTP(user)
	ctx.JSON(http.StatusOK, rsp)
}
