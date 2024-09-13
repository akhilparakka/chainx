package network

import (
	"bytes"
	"fmt"
	"sync"
)

type LocalTransport struct {
	addr NetAddr

	lock        sync.RWMutex
	peers       map[NetAddr]*LocalTransport
	consumeChan chan RPC
}

func NewLocalTransport(addr NetAddr) Transport {
	return &LocalTransport{
		addr:        addr,
		peers:       make(map[NetAddr]*LocalTransport),
		consumeChan: make(chan RPC, 1024),
	}
}

func (t *LocalTransport) Connect(tr Transport) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.peers[tr.Addr()] = tr.(*LocalTransport)

	return nil
}

func (t *LocalTransport) SendMessage(to NetAddr, msg []byte) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	peer, ok := t.peers[to]
	if !ok {
		return fmt.Errorf("%s: could not send message to %s", t.addr, to)
	}

	peer.consumeChan <- RPC{
		From:    t.addr,
		Payload: bytes.NewReader(msg),
	}

	return nil
}

func (t *LocalTransport) Consume() <-chan RPC {
	return t.consumeChan
}

func (t *LocalTransport) Addr() NetAddr {
	return t.addr
}
