package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) SetRoutes() {
	router := gin.Default()

	router.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Hello world")
	})

	s.Router = router
}
