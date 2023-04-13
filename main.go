package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/rohanchauhan02/simplebank/api"
	db "github.com/rohanchauhan02/simplebank/db/sqlc"
	"github.com/rohanchauhan02/simplebank/util"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can not load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("canot connect to db:", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("canot start server:", err)
	}
}
