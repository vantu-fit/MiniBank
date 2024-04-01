package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/vantu-fit/master-go-be/api"
	db "github.com/vantu-fit/master-go-be/db/sqlc"
	"github.com/vantu-fit/master-go-be/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot read env: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.BDSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	
	server , err  := api.NewServer(store)
	if err != nil {
		log.Fatal("cannot create server ")
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot connect to server :", err)

	}

}
