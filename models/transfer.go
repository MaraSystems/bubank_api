package models

import (
	db "github.com/MaraSystems/graybank_api/db/sqlc"
	"github.com/MaraSystems/graybank_api/utils"
)

type ListTransfersRequest struct {
	utils.PageRequest
	AccountID int64 `form:"account_id" binding:"required,gt=1"`
}

type CreateTransferRequest struct {
	Amount        int64 `json:"amount" binding:"required,gt=0"`
	FromAccountID int64 `json:"from_account_id" binding:"required,gt=1"`
	ToAccountID   int64 `json:"to_account_id" binding:"required,gt=1"`
}

type CreateTransferResponse struct {
	Transfer    db.Transfer `json:"transfer"`
	FromAccount db.Account  `json:"from_account"`
	ToAccount   db.Account  `json:"to_account"`
	FromEntry   db.Entry    `json:"from_entry"`
	ToEntry     db.Entry    `json:"to_entry"`
}
