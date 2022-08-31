package main

import (
	"context"
	"database/sql"
	"junk/boop-server/pgdb"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func run() error {
	ctx := context.Background()

	db, err := sql.Open("postgres", "dbname=boop sslmode=disable")
	if err != nil {
		return err
	}

	queries := pgdb.New(db)

	http.HandleFunc("/beans", MakeHttpHandler(handleBeans, ctx, queries))
	http.HandleFunc("/events", MakeHttpHandler(handleEvents, ctx, queries))
	http.HandleFunc("/tasks", MakeHttpHandler(handleTasks, ctx, queries))
	http.HandleFunc("/", MakeHttpHandler(handleRoot, ctx, queries))

	log.Print("Running on :22022")

	return http.ListenAndServe(":22022", nil)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
