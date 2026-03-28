package entries

import (
	"database/sql"
	"net/http"

	db "github.com/MaraSystems/graybank_api/db/sqlc"
	"github.com/MaraSystems/graybank_api/domains/accounts"
	"github.com/MaraSystems/graybank_api/utils"
	"github.com/gin-gonic/gin"
)

func (h *EntryHandler) createEntry(ctx *gin.Context) {
	var req CreateEntryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	_, err := accounts.GetAccountService(ctx, h.server.Store, req.AccountID, "account")
	if err != nil {
		return
	}

	arg := db.UpdateAccountBalanceParams{
		Amount: req.Amount,
		ID:     req.AccountID,
	}

	res := CreateEntryResponse{}
	err = h.server.Store.ExecuteTx(ctx.Request.Context(), func(q *db.Queries) error {
		var err error
		res.Account, err = q.UpdateAccountBalance(ctx, arg)
		if err != nil {
			return err
		}

		res.Entry, err = q.CreateEntry(ctx, db.CreateEntryParams{
			AccountID: sql.NullInt64{Int64: req.AccountID, Valid: true},
			Amount:    req.Amount,
			Model:     "Topup",
		})

		return err
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, res)
}
