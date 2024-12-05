package main

import (
	"gopher_social/internal/db"
	"gopher_social/internal/store"
	"log"
)

func main() {
	dbAddr := "postgres://admin:adminpassword@localhost:5432/gopher_social?sslmode=disable"
	conn, err := db.New(dbAddr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	store := store.NewPostgresStorage(conn)
	db.Seed(store, conn)
}
