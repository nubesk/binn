package binn

import (
	"time"
)

type Container interface {
	Message() (*Message)
	ID() string
	ExpiredAt() *time.Time
}

type Message struct {
	Text string
}

type Bottle struct {
	id        string
	message   *Message
	expiredAt *time.Time
}

func NewBottle(id string, text string, expiredAt *time.Time) *Bottle {
	return &Bottle{
		id:      id,
		message: &Message{
			Text: text,
		},
		expiredAt: expiredAt,
	}
}

func (b *Bottle) Message() *Message {
	return b.message
}

func (b *Bottle) ID() string {
	return b.id
}

func (b *Bottle) ExpiredAt() *time.Time {
	return b.expiredAt
}
