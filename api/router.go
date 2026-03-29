package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingExample godoc
//
//	@Summary	Ping API
//	@Schemes
//	@Description	do ping
//	@Tags			Ping
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	Hello	world
//	@Router			/ [get]
func (s *Server) ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "Hello world")
}

func (s *Server) SetRoutes() {
	s.Router = gin.Default()
	s.Router.GET("", s.ping)
}
