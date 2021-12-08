package main

import (
	"sync"
)

type ChatData []Message

type Chat struct {
	Data ChatData
	Mu   sync.RWMutex
}

type ChatStorage interface {
	SendToChat(*Message) *Message
	GetChat(int) []Message
}

func newChatStorage() *Chat {
	return &Chat{}
}

func (c *Chat) SendToChat(send *Message) *Message {
	c.Mu.Lock()
	c.Data = append(c.Data, *send)
	c.Mu.Unlock()

	return send
}

func (c *Chat) GetChat(page int) []Message {
	c.Mu.RLock()
	v := getMessageByPage(c.Data, page)
	c.Mu.RUnlock()

	return v
}
