package eventsource

type Eventstore struct {
	id        string   `json:"ID"`
	version   int      `json:"Version"`
	data      struct{} `json:"Data"`
	timestamp string   `json:"Time"`
}

type data struct {
	name     string `json:"Name`
	age      int    `json:"Age`
	marriage bool   `json:"Status"`
}
