package auth

import (
	"fmt"
	"net/http"

	"github.com/MaraSystems/graybank_api/domains/users"
	"github.com/MaraSystems/graybank_api/models"
	"github.com/MaraSystems/graybank_api/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// @Summary		User login
// @Description	Log in user
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			request	body		models.LoginParams	true	"credentials of the user"
// @Success		201		{object}	models.LoginResponseParams
// @Router			/auth [post]
func (h AuthHandler) loginUser(ctx *gin.Context) {
	var req models.LoginParams
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	user, err := h.server.Store.GetUser(ctx.Request.Context(), req.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(fmt.Errorf("user not found: %s", req.Username)))
			return
		}

		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	err = utils.VerifyPassword(user.HashedPassword, req.Password)
	if err != nil {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse(fmt.Errorf("password is incorrect")))
		return
	}

	token, err := h.server.TokenMaker.Create(req.Username, h.server.Config.AccessDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	res := models.LoginResponseParams{
		Token: token,
		User:  users.UserResponse(user),
	}

	ctx.JSON(http.StatusCreated, res)
}
