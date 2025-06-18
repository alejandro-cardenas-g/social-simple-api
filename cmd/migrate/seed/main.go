package main

import (
	"log"

	"github.com/alejandro-cardenas-g/social/internal/db"
	"github.com/alejandro-cardenas-g/social/internal/env"
	"github.com/alejandro-cardenas-g/social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:password@localhost/socialnetwork?sslmode=disable")
	log.Println(addr)
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)
	db.Seed(store, conn)
}
