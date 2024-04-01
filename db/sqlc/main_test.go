package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/vantu-fit/master-go-be/utils"
)

var testQueries *Queries
var testDb *sql.DB
var store Store

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")
	testDb, err := sql.Open(config.DBDriver, config.BDSource)
	fmt.Println(config.BDSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDb)
	store = NewStore(testDb)
	os.Exit(m.Run())
}
