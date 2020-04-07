package domain

// Message represent the structure of a single value
// eg:- key
//		value
type Message struct {
	Format  string      `json:"format,omitempty"` //avro,json,string,byte, default byte
	Subject string      `json:"subject,omitempty"`
	Version int         `json:"version,omitempty"`
	Body    interface{} `json:"body,omitempty"`
	Schema  interface{} `json:"schema,omitempty"`
}

//Event represent the structure of request message
type Event struct {
	Topic   string            `json:"topic"`
	Headers map[string]string `json:"headers"`
	Key     Message           `json:"key,omitempty"`
	Value   Message           `json:"value,omitempty"`
}
