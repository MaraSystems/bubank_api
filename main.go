package main

import (
	"database/sql"
	"log"

	"github.com/MaraSystems/graybank_api/api"
	db "github.com/MaraSystems/graybank_api/db/sqlc"
	"github.com/MaraSystems/graybank_api/domains/accounts"
	"github.com/MaraSystems/graybank_api/domains/auth"
	"github.com/MaraSystems/graybank_api/domains/entries"
	"github.com/MaraSystems/graybank_api/domains/transfers"
	"github.com/MaraSystems/graybank_api/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to the database: ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	AddDomainRoutes(server)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}

func AddDomainRoutes(server *api.Server) {
	auth.SetAuthRoutes(server)
	accounts.SetAccountsRoutes(server)
	entries.SetEntriesRoutes(server)
	transfers.SetTransfersRoutes(server)
}
