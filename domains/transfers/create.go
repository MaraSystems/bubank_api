package transfers

import (
	"database/sql"
	"net/http"

	db "github.com/MaraSystems/bubank_api/db/sqlc"
	"github.com/MaraSystems/bubank_api/domains/accounts"
	"github.com/MaraSystems/bubank_api/models"
	"github.com/MaraSystems/bubank_api/utils"
	"github.com/gin-gonic/gin"
)

func (h TransferHandler) createTransfer(ctx *gin.Context) {
	var req models.CreateTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	_, err := accounts.GetAccountService(ctx, h.server.Store, req.FromAccountID, "account")
	if err != nil {
		return
	}

	arg := db.CreateTransferParams{
		Amount:        req.Amount,
		FromAccountID: sql.NullInt64{Int64: req.FromAccountID, Valid: true},
		ToAccountID:   sql.NullInt64{Int64: req.ToAccountID, Valid: true},
	}

	res := models.CreateTransferResponse{}
	err = h.server.Store.ExecuteTx(ctx.Request.Context(), func(q *db.Queries) error {
		res.Transfer, err = q.CreateTransfer(ctx, arg)
		if err != nil {
			return err
		}

		res.FromAccount, err = q.UpdateAccountBalance(ctx, db.UpdateAccountBalanceParams{
			ID:     req.FromAccountID,
			Amount: -req.Amount,
		})
		if err != nil {
			return err
		}

		res.ToAccount, err = q.UpdateAccountBalance(ctx, db.UpdateAccountBalanceParams{
			ID:     req.ToAccountID,
			Amount: req.Amount,
		})
		if err != nil {
			return err
		}

		res.FromEntry, err = q.CreateEntry(ctx, db.CreateEntryParams{
			AccountID: sql.NullInt64{Int64: req.FromAccountID, Valid: true},
			Amount:    -req.Amount,
			Model:     "Transfer",
		})
		if err != nil {
			return err
		}

		res.ToEntry, err = q.CreateEntry(ctx, db.CreateEntryParams{
			AccountID: sql.NullInt64{Int64: req.ToAccountID, Valid: true},
			Amount:    req.Amount,
			Model:     "Transfer",
		})

		return err
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, res)
}
