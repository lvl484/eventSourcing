package eventsource

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	addr := flag.String("a", ":80", "address of app")
	flag.Parse()

	event := newEventStoreDb()
	defer event.db.Close()

	mainRoute := mux.NewRouter()
	apiRoute := mainRoute.PathPrefix("/api/v1").Subrouter()
	apiRoute.HandleFunc("/list", event.AddEvent).Methods("POST")
	apiRoute.HandleFunc("/list/{ID}", event.InsertEvent).Methods("PUT")
	apiRoute.HandleFunc("/list/{id, version}", event.GetEventByVersion).Methods("GET")
	apiRoute.HandleFunc("/list/{id, timestamp}", event.GetEventByTimestamp).Methods("GET")
	apiRoute.HandleFunc("/list/{id, versionMin, versionMax, PAGE}", event.GetEventByRangeOfVersion).Methods("GET")
	apiRoute.HandleFunc("/list/{id, timeMin, timeMax, PAGE}", event.GetEventByRangeOfTimestamp).Methods("GET")

	if err := http.ListenAndServe(*addr, mainRoute); err != nil {
		log.Fatal(err.Error())
	}
}
