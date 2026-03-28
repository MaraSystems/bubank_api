package accounts

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/MaraSystems/graybank_api/db/sqlc"
	"github.com/MaraSystems/graybank_api/middlewares"
	"github.com/MaraSystems/graybank_api/utils"
	"github.com/gin-gonic/gin"
)

func GetAccountService(ctx *gin.Context, store db.Store, id int64, entity string) (*db.Account, error) {
	account, err := store.GetAccount(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(fmt.Errorf("account not found: %d", id)))
			return nil, err
		}

		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return nil, err
	}

	authorization := ctx.MustGet(middlewares.AuthKey).(*utils.TokenPayload)
	if authorization.Username != account.Owner {
		err = fmt.Errorf("%s does not belong to the authenticated user", entity)
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return nil, err
	}

	return &account, nil
}

func (h *AccountHandler) getAccount(ctx *gin.Context) {
	var req utils.IDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	account, err := GetAccountService(ctx, h.server.Store, req.ID, "account")
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, account)
}
