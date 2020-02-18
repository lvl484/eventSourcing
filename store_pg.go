package eventsource
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

const offset = 10
const tempUUID = 11111111-1111-1111-1111-111111111111

type EventStoreDb struct {
	db := pg.Connect(pgOptions())
defer db.Close()
}

func newEventStoreDb() *EventStoreDb {
	
	database, err := sql.Open("postgres", "user=customer password=myorder dbname=eventstore host=localhost ")

	if err != nil {
		log.Println(err)
	}

	return &EventStoreDb{
		db: database,
	}
}

func (event *EventStoreDb) AddEvent(w http.ResponseWriter, r *http.Request) {
	var event EventStore

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	event.id = uuid.New().String()

	_, err = event.db.Exec(
		"CREATE TABLE ?
 (
	id uuid NOT NULL,
	version bigserial NOT NULL , 
	data json NOT NULL,
    time timestamp with time zone NOT NULL DEFAULT now(),
  CONSTRAINT event_pkey PRIMARY KEY (version)
);", event.id
	)

	if err != nil {
		log.Println(err)
		return
	}

	__, err = event.db.Exec("INSERT INTO ? (ID,Data) VALUES(?,?)",
		event.id, event.id, event.data)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (event *EventStoreDb) InsertEvent(w http.ResponseWriter, r *http.Request) {
	var event EventStore

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	
	__, err = event.db.Exec("INSERT INTO ? (ID,Data) VALUES(?,?)",
		event.id, event.id, event.data)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}




func (event *EventStoreDb) GetEventByVersion(w http.ResponseWriter, r *http.Request) {
	var ev EventStore

	params := mux.Vars(r)

	row := event.db.QueryRow("SELECT * FROM ? WHERE version=?",params["id"] params["version"])

	err := row.Scan(&ev.version,&ev.id, &ev.data, &ev.timestamp)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.NewEncoder(w).Encode(h)
	if err != nil {
		log.Println(err.Error())
	}
}

func (event *EventStoreDb) GetGetEventByTimestamp(w http.ResponseWriter, r *http.Request) {
	var ev EventStore

	params := mux.Vars(r)

	row := event.db.QueryRow("SELECT * FROM ? WHERE timestamp=?",params["id"] params["timestamp"])

	err := row.Scan(&ev.version,&ev.ID, &ev.data, &ev.timestamp)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.NewEncoder(w).Encode(h)
	if err != nil {
		log.Println(err.Error())
	}
}
func (event *EventStoreDb) GetEventByRangeOfVersion(w http.ResponseWriter, r *http.Request) {
	var events []EventStore

	params := mux.Vars(r)
	page, err := strconv.Atoi(params["PAGE"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	page--

	offsetnum := page * offset

	rows, err := event.db.Query("SELECT * FROM ? WHERE version BETWEEN ? AND ? LIMIT ?,?", 
	params["id"], params["versionMin"], param["versionMax"], offsetnum, offset)
	if err != nil {
		log.Println(err)
		return
	}

	for rows.Next() {
		var  ev EventStore
		err = rows.Scan(&ev.id, &ev.version, &ev.data, &h.timestamp)

		if err != nil {
			log.Println(err)
			continue
		}

		if err = rows.Err(); err != nil {
			log.Println(err)
			continue
		}

		event = append(event, h)
	}

	err = json.NewEncoder(w).Encode(event)
	if err != nil {
		log.Println(err)
	}
}
func (event *EventStoreDb) GetEventByRangeOfVersion(w http.ResponseWriter, r *http.Request) {
	var events []EventStore

	params := mux.Vars(r)
	page, err := strconv.Atoi(params["PAGE"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	page--

	offsetnum := page * offset

	rows, err := event.db.Query("SELECT * FROM ? WHERE version BETWEEN ? AND ? LIMIT ?,?", 
	params["id"], params["versionMin"], param["versionMax"], offsetnum, offset)
	if err != nil {
		log.Println(err)
		return
	}

	for rows.Next() {
		var  ev EventStore
		err = rows.Scan(&ev.id, &ev.version, &ev.data, &ev.timestamp)

		if err != nil {
			log.Println(err)
			continue
		}

		if err = rows.Err(); err != nil {
			log.Println(err)
			continue
		}

		event = append(event, h)
	}

	err = json.NewEncoder(w).Encode(event)
	if err != nil {
		log.Println(err)
	}
}

func (event *EventStoreDb) GetEventByRangeOfTimestamp(w http.ResponseWriter, r *http.Request) {
	var events []EventStore

	params := mux.Vars(r)
	page, err := strconv.Atoi(params["PAGE"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	page--

	offsetnum := page * offset

	rows, err := event.db.Query("SELECT * FROM ? WHERE time BETWEEN ? AND ? LIMIT ?,?", 
	params["id"], params["timeMin"], param["timeMax"], offsetnum, offset)
	if err != nil {
		log.Println(err)
		return
	}

	for rows.Next() {
		var  ev EventStore
		err = rows.Scan(&ev.id, &ev.version, &ev.data, &ev.timestamp)

		if err != nil {
			log.Println(err)
			continue
		}

		if err = rows.Err(); err != nil {
			log.Println(err)
			continue
		}

		event = append(event, h)
	}

	err = json.NewEncoder(w).Encode(event)
	if err != nil {
		log.Println(err)
	}
}










