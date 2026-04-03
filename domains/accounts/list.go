package accounts

import (
	"net/http"

	db "github.com/MaraSystems/bubank_api/db/sqlc"
	"github.com/MaraSystems/bubank_api/middlewares"
	"github.com/MaraSystems/bubank_api/models"
	"github.com/MaraSystems/bubank_api/utils"
	"github.com/gin-gonic/gin"
)

func (h *AccountHandler) listAccounts(ctx *gin.Context) {
	var req models.ListAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	authorization := ctx.MustGet(middlewares.AuthKey).(*utils.TokenPayload)

	limit := int32(10)
	offset := int32(0)

	if req.Limit != nil {
		limit = *req.Limit
	}
	if req.Offset != nil {
		offset = *req.Offset
	}

	arg := db.ListAccountsParams{
		Owner:  authorization.Username,
		Limit:  limit,
		Offset: offset * limit,
	}

	accounts, err := h.server.Store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
