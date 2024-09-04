package core

import (
	"bytes"
	"encoding/gob"

	"github.com/akhilparakka/chainx/types"
)

type Header struct {
	Version       uint32
	DataHash      types.Hash
	PrevblockHash types.Hash
	Timestamp     int64
	Height        uint32
	Nonce         uint64
}

type Block struct {
	*Header
	Transactions []Transaction
}

func Newblock(h *Header, txs []Transaction) *Block {
	return &Block{
		Header:       h,
		Transactions: txs,
	}
}

func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(h)

	return buf.Bytes()
}
