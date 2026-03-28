package accounts

import (
	"github.com/MaraSystems/graybank_api/api"
	"github.com/MaraSystems/graybank_api/middlewares"
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
