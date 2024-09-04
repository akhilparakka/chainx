package network

type netAddr string

type RPC struct {
	From    netAddr
	Payload []byte
}

type Transport interface {
	Connect(Transport) error
	SendMessage(netAddr, []byte) error
	Consume() <-chan RPC
	Addr() netAddr
}
