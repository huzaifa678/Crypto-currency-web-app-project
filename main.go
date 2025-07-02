package main

import (
	"context"
	"database/sql"
	"embed"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/huzaifa678/Crypto-currency-web-app-project/api"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	"github.com/huzaifa678/Crypto-currency-web-app-project/gapi"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//go:embed docs/*
var docsFS embed.FS

func main() {

	config, err := utils.LoadConfig(".")

	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	conn, err := sql.Open(config.Dbdriver, config.Dbsource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	/*server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal("failed to create server:", err)
	}
	
	err = server.Start(config.HTTPServerAddr)
	if err != nil {
		log.Fatal("failed to start the server:", err)
	}*/

	go runGatewayServer(config, store)
	runGrpcServer(config, store)
}

func runGrpcServer(config utils.Config, store db.Store_interface) {
	server, err := gapi.NewServer(store, config)

	if err != nil {
		log.Fatal("Cannot create the GRPC server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCryptoWebAppServer(grpcServer, server)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddr)
	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal("Cannot start the server:", err)
	}
}

func runGatewayServer(config utils.Config, store db.Store_interface) {
	server, err := gapi.NewServer(store, config)

	if err != nil {
		log.Fatal("Cannot create the GRPC server:", err)
	}

	grpcMux := runtime.NewServeMux()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterCryptoWebAppHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("Failed to register to the gateway handler server:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	fs := http.FileServer(http.FS(docsFS))
	mux.Handle("/docs/", http.StripPrefix("/", fs))

	listener, err := net.Listen("tcp", config.HTTPServerAddr)
	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("start HTTP server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)

	if err != nil {
		log.Fatal("Cannot start the HTTP gateway server:", err)
	}
}

func runGinServer(config utils.Config, store db.Store_interface) {
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal("failed to create server:", err)
	}
	
	err = server.Start(config.HTTPServerAddr)
	if err != nil {
		log.Fatal("failed to start the server:", err)
	}
}