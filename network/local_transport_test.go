package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalTransportConnect(t *testing.T) {
	tra := NewLocalTransport("A").(*LocalTransport)
	trb := NewLocalTransport("B").(*LocalTransport)

	tra.Connect(trb)
	trb.Connect(tra)

	assert.Equal(t, tra.peers[trb.Addr()], trb)
	assert.Equal(t, trb.peers[tra.Addr()], tra)

}

func TestSendMessage(t *testing.T) {
	tra := NewLocalTransport("A").(*LocalTransport)
	trb := NewLocalTransport("B").(*LocalTransport)

	msg := []byte("Hello A")

	tra.Connect(trb)
	assert.NotNil(t, trb.SendMessage(tra.Addr(), msg))

	trb.Connect(tra)
	assert.Nil(t, trb.SendMessage(tra.Addr(), msg))

	rpc := <-tra.Consume()

	assert.Equal(t, rpc.From, trb.Addr())
	assert.Equal(t, rpc.Payload, msg)

}
