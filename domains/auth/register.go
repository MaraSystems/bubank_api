package auth

import (
	"errors"
	"fmt"
	"net/http"

	db "github.com/MaraSystems/graybank_api/db/sqlc"
	"github.com/MaraSystems/graybank_api/domains/users"
	"github.com/MaraSystems/graybank_api/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

func (h AuthHandler) registerUser(ctx *gin.Context) {
	println("-->> Here")
	var req RegisterParams
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

	fmt.Println(arg)

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

	res := users.UserResponse(user)
	ctx.JSON(http.StatusCreated, res)
}
