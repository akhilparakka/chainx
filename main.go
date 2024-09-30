package main

import (
	"bytes"
	"fmt"
	"log"
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
	trremoteA := network.NewLocalTransport("REMOTE A")
	trremoteB := network.NewLocalTransport("REMOTE B")
	trremoteC := network.NewLocalTransport("REMOTE C")

	trlocal.Connect(trremoteA)
	trremoteA.Connect(trremoteB)
	trremoteB.Connect(trremoteC)

	trremoteA.Connect(trlocal)

	initRemoteServer([]network.Transport{trremoteA, trremoteB, trremoteC})

	go func() {
		for {
			if err := sendtransaction(trremoteA, trlocal.Addr()); err != nil {
				logrus.Error(err)
			}
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		trLate := network.NewLocalTransport("LATE_REMOTE")
		trremoteC.Connect(trLate)
		lateServer := makeServer(string(trLate.Addr()), trLate, nil)
		time.Sleep(7 * time.Second)
		go lateServer.Start()
	}()

	privKey := crypto.GeneratePrivateKey()

	localServer := makeServer("LOCAL", trlocal, &privKey)

	localServer.Start()
}

func initRemoteServer(trs []network.Transport) {
	for i := 0; i < len(trs); i++ {
		id := fmt.Sprintf("REMOTE_%d", i)
		s := makeServer(id, trs[i], nil)

		go s.Start()
	}
}

func makeServer(id string, tr network.Transport, privKey *crypto.PrivateKey) *network.Server {
	opts := network.Serveropts{
		PrivateKey: privKey,
		ID:         id,
		Transports: []network.Transport{tr},
	}

	s, err := network.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}

	return s
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
