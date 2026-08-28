package main

import (
	"aromaatelier/internal/flow031"
	"aromaatelier/internal/httpapi"
	"aromaatelier/internal/store"
	"flag"
	"log"
	"net/http"
)

func main() {
	path := flag.String("db", "aroma.db", "path to embedded database")
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()
	db, err := store.Open(*path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	server := httpapi.New(flow031.New(db))
	log.Printf("aroma atelier listening on %s", *addr)
	if err := http.ListenAndServe(*addr, server.Handler()); err != nil {
		log.Fatal(err)
	}
}
