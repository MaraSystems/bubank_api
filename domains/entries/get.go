package entries

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/MaraSystems/bubank_api/domains/accounts"
	"github.com/MaraSystems/bubank_api/utils"
	"github.com/gin-gonic/gin"
)

func (h EntryHandler) getEntry(ctx *gin.Context) {
	var req utils.IDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	entry, err := h.server.Store.GetEntry(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(fmt.Errorf("entry not found: %d", req.ID)))
			return
		}

		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	if _, err = accounts.GetAccountService(ctx, h.server.Store, entry.AccountID.Int64, "entry"); err != nil {
		return
	}

	ctx.JSON(http.StatusOK, entry)
}
