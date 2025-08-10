package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

var testStore Store_interface

func TestMain(m *testing.M) {

	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	
	connPool, err := pgxpool.New(context.Background(), config.Dbsource)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	
	testStore = NewStore(connPool)

	os.Exit(m.Run())
}


