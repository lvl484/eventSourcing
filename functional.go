package main

import (
	"strings"
	"time"
)

//Creare and format field ev.Timestamp
func Timestamp() string {
	timestamp := time.Now().String()
	lastBin := strings.LastIndex(timestamp, " +")
	return timestamp[0:lastBin]
}

//Store lastversion(field ev.Version) for each uuid(field ev.Id)
var LastVersionStore = make(map[string]int)

func Lastversion(uid string, ver int) {

	LastVersionStore["LastVer"] = ver
	LastVersionStore[uid] = ver
}
