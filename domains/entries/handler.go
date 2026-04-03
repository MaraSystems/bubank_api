package entries

import (
	"github.com/MaraSystems/bubank_api/api"
	"github.com/MaraSystems/bubank_api/middlewares"
)

type EntryHandler struct {
	server *api.Server
}

func SetEntriesRoutes(server *api.Server) {
	handler := &EntryHandler{
		server: server,
	}

	router := server.Router.Group("/entries")

	router.POST("", middlewares.AuthMiddleWare(server.TokenMaker), handler.createEntry)
	router.GET("", middlewares.AuthMiddleWare(server.TokenMaker), handler.liseEntries)
	router.GET(":id", middlewares.AuthMiddleWare(server.TokenMaker), handler.getEntry)
}
