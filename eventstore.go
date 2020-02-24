package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Eventstore struct {
	Version   int    `json:"version"`
	Id        string `json:"id"`
	Data      DATA   `json:"data"`
	Timestamp string `json:"time"`
}

type DATA map[string]interface{}

//Value method marshall the map to JSONB data (= []byte)
func (p DATA) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	return j, err
}

func (p *DATA) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*p, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("Type assertion .(map[string]interface{}) failed.")
	}

	return nil
}
