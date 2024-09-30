package network

import (
	"bytes"
	"os"
	"time"

	"github.com/akhilparakka/chainx/core"
	"github.com/akhilparakka/chainx/crypto"
	"github.com/akhilparakka/chainx/types"
	"github.com/go-kit/log"
)

var DefaultBlockTime = 5 * time.Second

type Serveropts struct {
	ID            string
	Logger        log.Logger
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	Transports    []Transport
	BlockTime     time.Duration
	PrivateKey    *crypto.PrivateKey
	blockchain    *core.Blockchain
}

type Server struct {
	Serveropts
	memPool     *TxPool
	chain       *core.Blockchain
	rpcChan     chan RPC
	endChan     chan struct{}
	isValidator bool
}

func NewServer(opts Serveropts) (*Server, error) {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = DefaultBlockTime
	}

	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}

	if opts.Logger == nil {
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		opts.Logger = log.With(opts.Logger, "ID", opts.ID)
	}

	chain, err := core.NewBlockchain(opts.Logger, genesisBlock())
	if err != nil {
		return nil, err
	}

	s := &Server{
		Serveropts:  opts,
		chain:       chain,
		memPool:     NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcChan:     make(chan RPC),
		endChan:     make(chan struct{}),
	}

	// If we dont have any processor from server options, use server as default.
	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}

	if s.isValidator {
		go s.validatorloop()
	}

	return s, nil

}

func (s *Server) Start() {
	s.initTransport()

free:
	for {
		select {
		case rpc := <-s.rpcChan:
			msg, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				s.Logger.Log("error", err)
			}

			if err := s.RPCProcessor.ProcessMessage(msg); err != nil {
				s.Logger.Log("error", err)

			}
		case <-s.endChan:
			break free

		}
	}

	s.Logger.Log("msg", "Server is shutting down")
}

func (s *Server) validatorloop() {
	ticker := time.NewTicker(s.BlockTime)

	s.Logger.Log("msg", "Starting validator loop", "blocktime", s.BlockTime)

	for {
		<-ticker.C
		s.createnewBlock()
	}
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
		return nil
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	tx.SetFirstSeen(time.Now().UnixNano())

	s.Logger.Log(
		"msg", "adding new tx to mempool",
		"hash", hash,
		"mempoolLength", s.memPool.Len(),
	)

	go s.broadcastTx(tx)

	return s.memPool.Add(tx)
}

func (s *Server) broadcastBlock(b *core.Block) error {
	// _ := &bytes.Buffer{}

	return nil
}

func (s *Server) broadcastTx(tx *core.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeTx, buf.Bytes())

	return s.broadcast(msg.Bytes())
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

func (s *Server) createnewBlock() error {
	currentHeader, err := s.chain.GetHeader(s.chain.Height())
	if err != nil {
		return err
	}

	// For now we are going to use all transactions that are in the mempool
	// Later on when we know the internal structure of the transaction,
	// Will implement some kind of complexity function which determins how
	// many transactions can be included in a block.
	txx := s.memPool.Transactions()

	block, err := core.NewBlockFromPrevHeader(currentHeader, txx)
	if err != nil {
		return err
	}

	block.Sign(*s.PrivateKey)

	if err := s.chain.AddBlock(block); err != nil {
		return err
	}

	s.memPool.Flush()

	return nil
}

func genesisBlock() *core.Block {
	header := &core.Header{
		Version:   1,
		DataHash:  types.Hash{},
		Height:    0,
		Timestamp: 000000,
	}

	b, _ := core.Newblock(header, nil)

	return b
}
