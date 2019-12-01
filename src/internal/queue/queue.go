package queue

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"time"

	"github.com/oklog/ulid"
	"github.com/streadway/amqp"
)

type (
	Config struct {
		Name      string
		Mandatory bool
		Immediate bool

		Durable    bool
		AutoDelete bool
		Exclusive  bool
		NoWait     bool
		Pass       string
		User       string
		Host       string
		Port       int
	}

	Queue struct {
		name      string
		mandatory bool
		immediate bool
		ch        *amqp.Channel
	}
)

func New(ch *amqp.Channel, cfg Config) *Queue {
	return &Queue{
		name:      cfg.Name,
		mandatory: cfg.Mandatory,
		immediate: cfg.Immediate,
		ch:        ch,
	}
}

const contentType = "application/json"

type Event struct {
	Extra    interface{} `json:"extra"`
	Callback string      `json:"callback"`
}

type Card struct {
	Card string `json:"card"`
	Event
}

type User struct {
	User string `json:"user"`
	Event
}

func (s *Queue) Send(event interface{}, appID string) error {
	js, err := json.Marshal(event)
	if err != nil {
		return wrapErr("json marshal: %w", err)
	}

	publishing := amqp.Publishing{
		ContentType: contentType,
		MessageId:   ulid.MustNew(ulid.Now(), rand.Reader).String(),
		Timestamp:   time.Now(),
		AppId:       appID,
		Body:        js,
	}

	return wrapErr("publish: %w", s.ch.Publish("", s.name, s.mandatory, s.immediate, publishing))
}

func wrapErr(str string, err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf(str, err)
}
