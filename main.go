package main

import (
	"time"

	"github.com/akhilparakka/chainx/network"
)

func main() {
	trlocal := network.NewLocalTransport("LOCAL")
	trremote := network.NewLocalTransport("REMOTE")

	trlocal.Connect(trremote)
	trremote.Connect(trlocal)

	go func() {
		for {
			trlocal.SendMessage(trremote.Addr(), []byte("Hi B!"))
			time.Sleep(1 * time.Second)
		}
	}()

	serveropts := network.Serveropts{
		Transports: []network.Transport{
			trremote,
		},
	}

	s := network.NewServer(serveropts)

	s.Start()
}
