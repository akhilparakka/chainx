package network

import (
	"io"
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
	buf := make([]byte, len(msg))
	n, err := rpc.Payload.Read(buf)
	assert.Nil(t, err)
	assert.Equal(t, n, len(msg))

	assert.Equal(t, rpc.From, trb.Addr())
	assert.Equal(t, buf, msg)
}

func TestBroadcast(t *testing.T) {
	tra := NewLocalTransport("A").(*LocalTransport)
	trb := NewLocalTransport("B").(*LocalTransport)
	trc := NewLocalTransport("C").(*LocalTransport)

	tra.Connect(trb)
	tra.Connect(trc)

	msg := []byte("helloo all")
	assert.Nil(t, tra.Broadcast(msg))
	rpcb := <-trb.Consume()
	b, err := io.ReadAll(rpcb.Payload)
	assert.Nil(t, err)
	assert.Equal(t, b, msg)

	rpcc := <-trc.Consume()
	c, err := io.ReadAll(rpcc.Payload)
	assert.Nil(t, err)
	assert.Equal(t, c, msg)
}
