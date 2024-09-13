package main

import (
	"bytes"
	"math/rand"
	"strconv"
	"time"

	"github.com/akhilparakka/chainx/core"
	"github.com/akhilparakka/chainx/crypto"
	"github.com/akhilparakka/chainx/network"
	"github.com/sirupsen/logrus"
)

func main() {
	trlocal := network.NewLocalTransport("LOCAL")
	trremote := network.NewLocalTransport("REMOTE")

	trlocal.Connect(trremote)
	trremote.Connect(trlocal)

	go func() {
		for {
			if err := sendtransaction(trremote, trlocal.Addr()); err != nil {
				logrus.Error(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	serveropts := network.Serveropts{
		Transports: []network.Transport{
			trlocal,
		},
	}

	s := network.NewServer(serveropts)

	s.Start()
}

func sendtransaction(tr network.Transport, to network.NetAddr) error {
	privkey := crypto.GeneratePrivateKey()
	data := []byte(strconv.FormatInt(int64(rand.Intn(1000)), 10))
	tx := core.NewTransaction(data)
	tx.Sign(privkey)
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}
	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())

	return tr.SendMessage(to, msg.Bytes())

}
