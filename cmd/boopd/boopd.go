package main

import (
	"context"
	"database/sql"
	"junk/boop-server"
	"junk/boop-server/pgdb"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func run() error {
	ctx := context.Background()

	db, err := sql.Open("postgres", "user=postgres host=192.168.42.89 dbname=boop sslmode=disable")
	if err != nil {
		return err
	}

	queries := pgdb.New(db)

	http.HandleFunc("/beans", boop.MakeHttpHandler(boop.HandleBeans, ctx, queries))
	http.HandleFunc("/events", boop.MakeHttpHandler(boop.HandleEvents, ctx, queries))
	http.HandleFunc("/tasks", boop.MakeHttpHandler(boop.HandleTasks, ctx, queries))
	http.HandleFunc("/", boop.MakeHttpHandler(boop.HandleRoot, ctx, queries))

	log.Print("Running on :22022")

	return http.ListenAndServe(":22022", nil)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
