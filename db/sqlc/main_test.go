package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vantu-fit/master-go-be/utils"
)

var testQueries *Queries
var store Store

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")
	testDb, err := pgxpool.New(context.Background(), config.BDSource)
	fmt.Println(config.BDSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(testDb)
	store = NewStore(testDb)

	defer testDb.Close()

	os.Exit(m.Run())
}
