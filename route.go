package main

import (
	"github.com/gorilla/mux"
)

func (event *EventStoreDb) NewRoute() *mux.Router {
	mainRoute := mux.NewRouter()
	mainRoute.Use(middleware)
	apiRoute := mainRoute.PathPrefix("/api/v1").Subrouter()
	apiRoute.HandleFunc("/list", event.InsertEventWithoutID).Methods("POST")
	apiRoute.HandleFunc("/list/{ID}", event.InsertEvent).Methods("POST")
	apiRoute.HandleFunc("/list/version/{VERSION}", event.GetEventByVersion).Methods("GET")
	apiRoute.HandleFunc("/list/{ID}", event.GetEventByRangeOfVersion).Methods("GET")
	apiRoute.HandleFunc("/list", event.GetEventByTimestamp).Methods("GET")
	apiRoute.HandleFunc("/list/time/{ID}", event.GetEventByRangeOfTimestamp).Methods("GET")

	return mainRoute
}
