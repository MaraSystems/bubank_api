package gapi

import (
	"testing"
	"time"

	db "github.com/MaraSystems/bubank_api/db/sqlc"
	pb "github.com/MaraSystems/bubank_api/pb"
	"github.com/MaraSystems/bubank_api/utils"
	"github.com/stretchr/testify/require"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Server struct {
	pb.UnimplementedBubankServer
	Store      db.Store
	Config     utils.Config
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
