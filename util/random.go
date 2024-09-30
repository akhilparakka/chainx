package util

import (
	"crypto/rand"
	"testing"

	"github.com/akhilparakka/chainx/core"
	"github.com/akhilparakka/chainx/crypto"
	"github.com/akhilparakka/chainx/types"
)

func RandomBytes(size int) []byte {
	token := make([]byte, size)
	rand.Read(token)
	return token
}

func RandomHash() types.Hash {
	return types.HashFromBytes(RandomBytes(32))
}

func NewRandomTransaction(size int) *core.Transaction {
	return core.NewTransaction(RandomBytes(size))
}

func NewRandomTransactionWithSignature(t *testing.T, privKey crypto.PrivateKey, size int) *core.Transaction {
	tx := NewRandomTransaction(size)
	tx.Sign(privKey)
	return tx
}
