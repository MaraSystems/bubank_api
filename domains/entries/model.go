package entries

import (
	db "github.com/MaraSystems/graybank_api/db/sqlc"
	"github.com/MaraSystems/graybank_api/utils"
)

type ListEntriesRequest struct {
	utils.PageRequest
	AccountID int64 `form:"account_id" binding:"required,gt=1"`
}

type CreateEntryRequest struct {
	Amount    int64 `json:"amount" binding:"required,gt=0"`
	AccountID int64 `json:"account_id" binding:"required,gt=1"`
}

type CreateEntryResponse struct {
	Account db.Account `json:"account"`
	Entry   db.Entry   `json:"entry"`
}
