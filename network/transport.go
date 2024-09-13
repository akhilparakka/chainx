package network

type NetAddr string

type Transport interface {
	Connect(Transport) error
	SendMessage(NetAddr, []byte) error
	Consume() <-chan RPC
	Addr() NetAddr
}
