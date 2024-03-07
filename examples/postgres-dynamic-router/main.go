package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bmviniciuss/forger-golang/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dsn := "postgres://forger_user:1234@localhost:5432/forger?sslmode=disable"
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return err
	}

	loader := NewPostgresLoader(db.DB)
	mux := mux.NewDynamicRouter(loader)

	fmt.Println("Server started at http://localhost:3000")
	http.ListenAndServe(":3000", mux)

	return nil
}
