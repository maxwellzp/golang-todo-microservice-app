package nats

import (
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

type Publisher struct {
	conn *nats.Conn
}

func NewPublisher() (*Publisher, error) {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://nats:4222"
	}
	conn, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}
	return &Publisher{conn: conn}, nil
}

func (p *Publisher) Publish(subject string, message string) {
	err := p.conn.Publish(subject, []byte(message))
	if err != nil {
		log.Printf("❌ Failed to publish to subject %s: %v\n", subject, err)
	}
}

func (p *Publisher) Close() {
	p.conn.Close()
}
