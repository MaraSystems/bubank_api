package auth

import (
	"github.com/MaraSystems/graybank_api/api"
	"github.com/MaraSystems/graybank_api/middlewares"
)

type AuthHandler struct {
	server *api.Server
}

func SetAuthRoutes(server *api.Server) {
	handler := &AuthHandler{
		server: server,
	}

	router := server.Router.Group("/auth")

	router.POST("/register", handler.registerUser)
	router.POST("", handler.loginUser)
	router.GET("", middlewares.AuthMiddleWare(server.TokenMaker), handler.getProfile)
}
