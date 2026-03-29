package main

import (
	"database/sql"
	"log"

	"github.com/MaraSystems/graybank_api/api"
	db "github.com/MaraSystems/graybank_api/db/sqlc"
	"github.com/MaraSystems/graybank_api/docs"
	"github.com/MaraSystems/graybank_api/domains/accounts"
	"github.com/MaraSystems/graybank_api/domains/auth"
	"github.com/MaraSystems/graybank_api/domains/entries"
	"github.com/MaraSystems/graybank_api/domains/transfers"
	"github.com/MaraSystems/graybank_api/utils"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Bank API
//	@version		1.0
//	@description	This is a dummy bank

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
	docs.SwaggerInfo.BasePath = ""
	server.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
