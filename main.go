package main

import (
	"database/sql"
	"github.com/SemmiDev/chi-bank/common"
	"log"

	"github.com/SemmiDev/chi-bank/api"
	db "github.com/SemmiDev/chi-bank/db/sqlc"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	c, err := common.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(c.DBDriver, c.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(c, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	if err := server.Start(c.ServerAddress); err != nil {
		log.Fatal(err.Error())
	}
}
