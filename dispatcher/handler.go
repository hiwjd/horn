package dispatcher

import (
	"log"
	"os"

	"github.com/nsqio/go-nsq"
)

type Handler struct {
	totalMessages int
	messagesShown int
}

func (th *Handler) HandleMessage(m *nsq.Message) error {
	th.messagesShown++
	_, err := os.Stdout.Write(m.Body)
	if err != nil {
		log.Fatalf("ERROR: failed to write to os.Stdout - %s", err)
	}
	_, err = os.Stdout.WriteString("\n")
	if err != nil {
		log.Fatalf("ERROR: failed to write to os.Stdout - %s", err)
	}
	if th.totalMessages > 0 && th.messagesShown >= th.totalMessages {
		os.Exit(0)
	}
	return nil
}

func NewHandler(total int) *Handler {
	return &Handler{totalMessages: total}
}
