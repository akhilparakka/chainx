package network

import (
	"fmt"
	"time"

	"github.com/akhilparakka/chainx/core"
	"github.com/akhilparakka/chainx/crypto"
	"github.com/sirupsen/logrus"
)

type Serveropts struct {
	Transports []Transport
	BlockTime  time.Duration
	PrivateKey *crypto.PrivateKey
}

type Server struct {
	Serveropts
	blockTime   time.Duration
	memPool     *TxPool
	rpcChan     chan RPC
	endChan     chan struct{}
	isValidator bool
}

func NewServer(opts Serveropts) *Server {
	return &Server{
		Serveropts:  opts,
		blockTime:   opts.BlockTime,
		memPool:     NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcChan:     make(chan RPC),
		endChan:     make(chan struct{}),
	}
}

func (s *Server) Start() {
	s.initTransport()
	ticker := time.NewTicker(s.blockTime)

free:
	for {
		select {
		case rpc := <-s.rpcChan:
			fmt.Printf("%+v\n", rpc)
		case <-s.endChan:
			break free
		case <-ticker.C:
			s.createnewBlock()
		}
	}

	fmt.Println("Server Shutdown")
}

func (s *Server) handleTransaction(tx *core.Transaction) error {
	if err := tx.Verify(); err != nil {
		return err
	}

	hash := tx.Hash(core.TxHasher{})

	if s.memPool.Has(hash) {
		logrus.WithFields(logrus.Fields{
			"hash": hash,
		}).Info("mempool already contains the transaction")
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"hash": hash,
	}).Info("adding new tx to the mempool")

	return s.memPool.Add(tx)
}

func (s *Server) createnewBlock() error {
	fmt.Println("Creating new block")
	return nil
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
