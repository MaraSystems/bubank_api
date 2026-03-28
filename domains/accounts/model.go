package accounts

import (
	db "github.com/MaraSystems/graybank_api/db/sqlc"
	"github.com/MaraSystems/graybank_api/utils"
)

type CreateAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

type ListAccountsRequest struct {
	utils.PageRequest
}

func DummyAccount() db.Account {
	return db.Account{
		Owner:    utils.RandomUsername(),
		Balance:  100,
		Currency: "NGN",
	}
}
