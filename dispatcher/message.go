package dispatcher

import "encoding/json"

type Chat struct {
	Id string
}

type Message struct {
	Type string
	Chat *Chat
	Data json.RawMessage
}

type MessageText struct {
	Text string
}

type MessageImage struct {
	Src    string
	Width  int
	Height int
	Size   int
}
