package accounts

import (
	"github.com/MaraSystems/bubank_api/api"
	db "github.com/MaraSystems/bubank_api/db/sqlc"
	"github.com/MaraSystems/bubank_api/middlewares"
	"github.com/MaraSystems/bubank_api/utils"
)

type AccountHandler struct {
	server *api.Server
}

func SetAccountsRoutes(server *api.Server) {
	handler := &AccountHandler{
		server: server,
	}

	router := server.Router.Group("/accounts")

	router.POST("", middlewares.AuthMiddleWare(server.TokenMaker), handler.createAccount)
	router.GET(":id", middlewares.AuthMiddleWare(server.TokenMaker), handler.getAccount)
	router.GET("", middlewares.AuthMiddleWare(server.TokenMaker), handler.listAccounts)
}

func dummyAccount() db.Account {
	return db.Account{
		Owner:    utils.RandomUsername(),
		Balance:  100,
		Currency: "NGN",
	}
}
