package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/MaraSystems/bubank_api/api"
	db "github.com/MaraSystems/bubank_api/db/sqlc"
	"github.com/MaraSystems/bubank_api/docs"
	"github.com/MaraSystems/bubank_api/domains/accounts"
	"github.com/MaraSystems/bubank_api/domains/auth"
	"github.com/MaraSystems/bubank_api/domains/entries"
	"github.com/MaraSystems/bubank_api/domains/transfers"
	"github.com/MaraSystems/bubank_api/gapi"
	"github.com/MaraSystems/bubank_api/pb"
	"github.com/MaraSystems/bubank_api/utils"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	go runGatewayServer(config, store)
	runGrpcServer(config, store)
}

func runGrpcServer(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBubankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener: ", err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start grpc server: ", err)
	}
}

func runGatewayServer(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	grpcMux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pb.RegisterBubankHandlerServer(ctx, grpcMux, server)

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create listener: ", err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start grpc server: ", err)
	}
}

func runGinServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	AddDomainRoutes(server)
	docs.SwaggerInfo.BasePath = ""
	server.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err = server.Start(config.HTTPServerAddress)
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
