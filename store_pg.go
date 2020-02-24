package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const insertEvent string = "INSERT INTO event (id, data,time) VALUES($1,$2,$3) RETURNING version;"
const selectVerEvent string = "SELECT  version, id, data, time FROM event WHERE version=$1;"
const selectTimeEvent string = "SELECT  version, id, data, time  FROM event WHERE time=$1;"
const selectRangVerEvent string = "SELECT version, id, data, time  FROM event WHERE (id=$1) AND (version>=$2) AND (version<=$3);"
const selectRangTimeEvent string = "SELECT version, id, data, time FROM event WHERE (id=$1) AND (time>=$2) AND (time<=$3);"

func (event *EventStoreDb) InsertEventWithoutID(w http.ResponseWriter, r *http.Request) {
	var ev Eventstore

	err := json.NewDecoder(r.Body).Decode(&ev.Data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	ev.Id = uuid.New().String()

	e, err := ev.Data.Value()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ev.Timestamp = Timestamp()

	err = event.db.QueryRow(insertEvent, ev.Id, e, ev.Timestamp).Scan(&ev.Version)
	if err != nil {
		log.Println(err)
		return
	}

	Lastversion(ev.Id, ev.Version)

	w.Write([]byte(ev.Id))
}
func (event *EventStoreDb) InsertEvent(w http.ResponseWriter, r *http.Request) {
	var ev Eventstore

	err := json.NewDecoder(r.Body).Decode(&ev.Data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	params := mux.Vars(r)

	_, err = uuid.Parse(params["ID"])
	if err == nil {
		ev.Id = params["ID"]

	} else {

		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	e, err := ev.Data.Value()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ev.Timestamp = Timestamp()
	err = event.db.QueryRow(insertEvent, ev.Id, e, ev.Timestamp).Scan(&ev.Version)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Lastversion(ev.Id, ev.Version)
	w.WriteHeader(http.StatusOK)
}

func (event *EventStoreDb) GetEventByVersion(w http.ResponseWriter, r *http.Request) {
	var ev Eventstore

	params := mux.Vars(r)
	ver, err := strconv.Atoi(params["VERSION"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	row := event.db.QueryRow(selectVerEvent, ver)

	err = row.Scan(&ev.Version, &ev.Id, &ev.Data, &ev.Timestamp)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = json.NewEncoder(w).Encode(ev)
	if err != nil {
		log.Println(err.Error())
	}
}
func (event *EventStoreDb) GetEventByRangeOfVersion(w http.ResponseWriter, r *http.Request) {
	var events []Eventstore

	params := mux.Vars(r)

	versionMin := (r.FormValue("vm"))
	if versionMin == "0" {
		versionMin = "1"
	}

	versionMax := r.FormValue("vM")
	if versionMax == "0" {
		versionMax = string(LastVersionStore[params["ID"]])
	}

	rows, err := event.db.Query(selectRangVerEvent, params["ID"], versionMin, versionMax)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var ev Eventstore
		err = rows.Scan(&ev.Version, &ev.Id, &ev.Data, &ev.Timestamp)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		events = append(events, ev)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		log.Println(err)
	}
}

func (event *EventStoreDb) GetEventByTimestamp(w http.ResponseWriter, r *http.Request) {
	var ev Eventstore
	var time = ""

	err := json.NewDecoder(r.Body).Decode(&time)
	if err != nil {
		log.Println(err)

		return
	}

	row := event.db.QueryRow(selectTimeEvent, &time)

	err = row.Scan(&ev.Version, &ev.Id, &ev.Data, &ev.Timestamp)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(ev)
	if err != nil {
		log.Println(err.Error())
	}
}

func (event *EventStoreDb) GetEventByRangeOfTimestamp(w http.ResponseWriter, r *http.Request) {

	var events []Eventstore

	type Times struct {
		TimesMin string `json:"timemin"`
		TimesMax string `json:"timemax"`
	}
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var t Times
	err = json.Unmarshal(b, &t)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	params := mux.Vars(r)

	rows, err := event.db.Query(selectRangTimeEvent, params["ID"], t.TimesMin, t.TimesMax)
	if err != nil {
		log.Println(err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var ev Eventstore
		err = rows.Scan(&ev.Version, &ev.Id, &ev.Data, &ev.Timestamp)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		events = append(events, ev)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		log.Println(err)
	}
}
