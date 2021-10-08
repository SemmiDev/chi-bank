package main

import (
	"database/sql"
	"log"

	"github.com/SemmiDev/chi-bank/api"
	db "github.com/SemmiDev/chi-bank/db/sqlc"
	"github.com/SemmiDev/chi-bank/util"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal(err.Error())
	}
}
