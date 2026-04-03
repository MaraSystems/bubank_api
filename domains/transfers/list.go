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

func (h *TransferHandler) listTransfers(ctx *gin.Context) {
	var req models.ListTransfersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	if _, err := accounts.GetAccountService(ctx, h.server.Store, req.AccountID, "account"); err != nil {
		return
	}

	limit := int32(10)
	offset := int32(0)

	if req.Limit != nil {
		limit = *req.Limit
	}
	if req.Offset != nil {
		offset = *req.Offset
	}

	arg := db.ListTransfersParams{
		FromAccountID: sql.NullInt64{Int64: req.AccountID, Valid: true},
		Limit:         limit,
		Offset:        offset * limit,
	}

	accounts, err := h.server.Store.ListTransfers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
