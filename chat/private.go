package main

import (
	"sync"
)

type PrivateData map[string][]Message

type Private struct {
	Data PrivateData
	Mu   sync.RWMutex
}

type PrivateStorage interface {
	SendToPrivate(string, *Message) *Message
	GetPrivate(string, int) []Message
}

func newPrivateStorage() *Private {
	return &Private{
		Data: make(PrivateData),
	}
}

func (p *Private) SendToPrivate(reciever string, send *Message) *Message {
	p.Mu.Lock()
	p.Data[reciever] = append(p.Data[reciever], *send)
	p.Mu.Unlock()

	return send
}

func (p *Private) GetPrivate(login string, page int) []Message {
	p.Mu.RLock()
	v := getMessageByPage(p.Data[login], page)
	p.Mu.RUnlock()

	return v
}
