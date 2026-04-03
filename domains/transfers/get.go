package transfers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/MaraSystems/bubank_api/domains/accounts"
	"github.com/MaraSystems/bubank_api/utils"
	"github.com/gin-gonic/gin"
)

func (h TransferHandler) getTransfer(ctx *gin.Context) {
	var req utils.IDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	transfer, err := h.server.Store.GetTransfer(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(fmt.Errorf("transfer not found: %d", req.ID)))
			return
		}

		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	if _, err = accounts.GetAccountService(ctx, h.server.Store, transfer.FromAccountID.Int64, "transfer"); err != nil {
		return
	}

	ctx.JSON(http.StatusOK, transfer)
}
