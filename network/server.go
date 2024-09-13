package network

import (
	"fmt"
	"time"

	"github.com/akhilparakka/chainx/core"
	"github.com/akhilparakka/chainx/crypto"
	"github.com/sirupsen/logrus"
)

var DefaultBlockTime = 5 * time.Second

type Serveropts struct {
	RPCHandler RPCHandler
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
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = DefaultBlockTime
	}

	s := &Server{
		Serveropts:  opts,
		blockTime:   opts.BlockTime,
		memPool:     NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcChan:     make(chan RPC),
		endChan:     make(chan struct{}),
	}

	if opts.RPCHandler == nil {
		opts.RPCHandler = NewDefaultRPCHandler(s)
	}

	s.Serveropts = opts

	return s

}

func (s *Server) Start() {
	s.initTransport()
	ticker := time.NewTicker(s.blockTime)

free:
	for {
		select {
		case rpc := <-s.rpcChan:
			if err := s.RPCHandler.HandleRPC(rpc); err != nil {
				logrus.Error(err)
			}
		case <-s.endChan:
			break free
		case <-ticker.C:
			s.createnewBlock()
		}
	}

	fmt.Println("Server Shutdown")
}

func (s *Server) ProcessTransaction(from NetAddr, tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})

	if s.memPool.Has(hash) {
		logrus.WithFields(logrus.Fields{
			"hash": hash,
		}).Info("mempool already contains the transaction")
		return nil
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	tx.SetFirstSeen(time.Now().UnixNano())

	logrus.WithFields(logrus.Fields{
		"hash":           hash,
		"mempool length": s.memPool.Len(),
	}).Info("adding new tx to the mempool")

	// TODO(): broadcast this tx to peers

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
