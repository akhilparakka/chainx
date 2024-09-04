package network

import "fmt"

type Serveropts struct {
	Transports []Transport
}

type Server struct {
	Serveropts
	rpcChan chan RPC
	endChan chan struct{}
}

func NewServer(opts Serveropts) *Server {
	return &Server{
		Serveropts: opts,
		rpcChan:    make(chan RPC),
		endChan:    make(chan struct{}),
	}
}

func (s *Server) Start() {
	s.initTransport()

free:
	for {
		select {
		case rpc := <-s.rpcChan:
			fmt.Printf("%+v\n", rpc)
		case <-s.endChan:
			break free

		}
	}
}

func (s *Server) initTransport() {
	for _, peer := range s.Transports {
		go func(peer Transport) {
			for rpc := range peer.Consume() {
				s.rpcChan <- rpc
			}
		}(peer)
	}
}
