package main

import (
	"github.com/huzaifa678/Crypto-currency-web-app-project/api"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	"database/sql"
	"log"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
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

	err = godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddr)
}