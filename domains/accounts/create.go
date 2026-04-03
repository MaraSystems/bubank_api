package accounts

import (
	"errors"
	"net/http"

	db "github.com/MaraSystems/bubank_api/db/sqlc"
	"github.com/MaraSystems/bubank_api/middlewares"
	"github.com/MaraSystems/bubank_api/models"
	"github.com/MaraSystems/bubank_api/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

func (h *AccountHandler) createAccount(ctx *gin.Context) {
	var req models.CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	authorization := ctx.MustGet(middlewares.AuthKey).(*utils.TokenPayload)

	arg := db.CreateAccountParams{
		Owner:    authorization.Username,
		Balance:  0,
		Currency: req.Currency,
	}

	account, err := h.server.Store.CreateAccount(ctx, arg)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			switch pgError.ConstraintName {
			case "accounts_owner_fkey", "owner_currency_key", "unique_violation":
				ctx.JSON(http.StatusForbidden, utils.ErrorResponse(err))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, account)
}
