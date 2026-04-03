package auth

import (
	"errors"
	"net/http"

	db "github.com/MaraSystems/bubank_api/db/sqlc"
	"github.com/MaraSystems/bubank_api/domains/users"
	"github.com/MaraSystems/bubank_api/models"
	"github.com/MaraSystems/bubank_api/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

// @Summary		User registration
// @Description	Create a new user
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			request	body		models.RegisterRequest	true	"username of the user"
// @Success		201		{object}	models.UserResponse
// @Router			/auth/register [post]
func (h AuthHandler) registerUser(ctx *gin.Context) {
	var req models.RegisterRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		Email:          req.Email,
		FullName:       req.FullName,
		HashedPassword: hashedPassword,
	}

	user, err := h.server.Store.CreateUser(ctx.Request.Context(), arg)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			switch pgError.Code {
			case "23505":
				ctx.JSON(http.StatusConflict, utils.ErrorResponse(errors.New("username or email is already in use")))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	res := users.UserToHTTP(user)
	ctx.JSON(http.StatusCreated, res)
}
