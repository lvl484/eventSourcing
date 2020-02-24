package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("a", ":8080", "address of app")
	flag.Parse()

	event := newEventStoreDb()
	defer event.db.Close()

	route := event.NewRoute()

	if err := http.ListenAndServe(*addr, route); err != nil {
		log.Fatal(err.Error())
	}
}
