package main

import (
	"crypto-system/api"
	db "crypto-system/db/sqlc"
	"database/sql"
	"log"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/crypto_db?sslmode=disable"
	serverAddr = "0.0.0.0:8081"
)
func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddr)
}