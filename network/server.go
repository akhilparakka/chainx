package network

import (
	"bytes"
	"fmt"
	"time"

	"github.com/akhilparakka/chainx/core"
	"github.com/akhilparakka/chainx/crypto"
	"github.com/sirupsen/logrus"
)

var DefaultBlockTime = 5 * time.Second

type Serveropts struct {
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	Transports    []Transport
	BlockTime     time.Duration
	PrivateKey    *crypto.PrivateKey
}

type Server struct {
	Serveropts
	memPool     *TxPool
	rpcChan     chan RPC
	endChan     chan struct{}
	isValidator bool
}

func NewServer(opts Serveropts) *Server {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = DefaultBlockTime
	}

	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}

	s := &Server{
		Serveropts:  opts,
		memPool:     NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcChan:     make(chan RPC),
		endChan:     make(chan struct{}),
	}

	// If we dont have any processor from server options, use server as default.
	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}

	return s

}

func (s *Server) Start() {
	s.initTransport()
	ticker := time.NewTicker(s.BlockTime)

free:
	for {
		select {
		case rpc := <-s.rpcChan:
			msg, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				logrus.Error(err)
			}

			if err := s.RPCProcessor.ProcessMessage(msg); err != nil {
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

func (s *Server) ProcessMessage(msg *DecodedMessage) error {

	switch t := msg.Data.(type) {
	case *core.Transaction:
		return s.processTransaction(t)
	}

	return nil
}

func (s *Server) broadcast(payload []byte) error {
	for _, tr := range s.Transports {
		if err := tr.Broadcast(payload); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) processTransaction(tx *core.Transaction) error {
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

	go s.broadcastTx(tx)

	return s.memPool.Add(tx)
}

func (s *Server) broadcastTx(tx *core.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeTx, buf.Bytes())

	return s.broadcast(msg.Bytes())
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
