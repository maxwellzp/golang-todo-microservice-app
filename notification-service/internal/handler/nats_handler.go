package handler

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

type NATSHandler struct {
	conn *nats.Conn
}

func NewNATSHandler(url string) (*NATSHandler, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NATSHandler{conn: conn}, nil
}

func (h *NATSHandler) Subscribe(subject string, handle func(msg string)) {
	_, err := h.conn.Subscribe(subject, func(m *nats.Msg) {
		handle(string(m.Data))
	})
	if err != nil {
		fmt.Printf("❌ Failed to subscribe to subject %s: %v\n", subject, err)
	} else {
		fmt.Printf("📡 Subscribed to subject: %s\n", subject)
	}
}
