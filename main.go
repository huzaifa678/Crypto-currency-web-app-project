package main

import (
	"database/sql"
	"log"

	api "github.com/huzaifa678/Crypto-currency-web-app-project/api"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	_ "github.com/lib/pq"
)


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
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal("failed to create server:", err)
	}
	
	err = server.Start(config.ServerAddr)
	if err != nil {
		log.Fatal("failed to start the server:", err)
	}

}