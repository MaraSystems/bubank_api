package models

import (
	"github.com/MaraSystems/graybank_api/utils"
)

type CreateAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

type ListAccountsRequest struct {
	utils.PageRequest
}
