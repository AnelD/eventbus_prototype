package bus

type Message struct {
	Type  string `json:"type"`
	Topic string `json:"topic"`
	Data  string `json:"data,omitempty"`
}
