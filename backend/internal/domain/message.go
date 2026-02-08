package domain

type Message struct {
	Type string `json:"type"`
	Data map[string]string
}
