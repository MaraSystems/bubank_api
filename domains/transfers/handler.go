package transfers

import (
	"github.com/MaraSystems/bubank_api/api"
	"github.com/MaraSystems/bubank_api/middlewares"
)

type TransferHandler struct {
	server *api.Server
}

func SetTransfersRoutes(server *api.Server) {
	handler := &TransferHandler{
		server: server,
	}

	router := server.Router.Group("/transfers")

	router.POST("", middlewares.AuthMiddleWare(server.TokenMaker), handler.createTransfer)
	router.GET("", middlewares.AuthMiddleWare(server.TokenMaker), handler.listTransfers)
	router.GET(":id", middlewares.AuthMiddleWare(server.TokenMaker), handler.getTransfer)
}
