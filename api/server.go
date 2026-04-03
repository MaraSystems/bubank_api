package api

import (
	"testing"
	"time"

	db "github.com/MaraSystems/bubank_api/db/sqlc"
	"github.com/MaraSystems/bubank_api/utils"
	"github.com/MaraSystems/bubank_api/validators"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Server struct {
	Store      db.Store
	Config     utils.Config
	Router     *gin.Engine
	TokenMaker *utils.TokenMaker
}

func NewServer(config utils.Config, store db.Store) (server *Server, err error) {
	tokenMaker, err := utils.NewTokenMaker(config.AccessSecretKey)
	if err != nil {
		return
	}

	server = &Server{
		Store:      store,
		Config:     config,
		TokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validators.ValidateCurrency)
	}

	server.SetRoutes()
	return
}

func TestServer(t *testing.T, store db.Store) (*Server, error) {
	config := utils.Config{
		AccessSecretKey: "12345678901234567890123456789012",
		AccessDuration:  time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server, nil
}

func (s *Server) Start(address string) error {
	return s.Router.Run(address)
}
