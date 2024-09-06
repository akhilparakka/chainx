package core

import (
	"fmt"

	"github.com/akhilparakka/chainx/crypto"
	"github.com/akhilparakka/chainx/types"
)

type Transaction struct {
	Data []byte

	From      crypto.PublicKey
	Signature *crypto.Signature

	hash types.Hash
}

func Newtransaction(data []byte) *Transaction {
	return &Transaction{
		Data: data,
	}
}

func (tx *Transaction) Hash(hasher Hasher[*Transaction]) types.Hash {
	if tx.hash.IsZero() {
		tx.hash = hasher.Hash(tx)
	}

	return hasher.Hash(tx)
}

func (tx *Transaction) Sign(privkey crypto.PrivateKey) {
	sig := privkey.Sign(tx.Data)

	tx.From = privkey.PublicKey()
	tx.Signature = sig
}

func (tx *Transaction) Verify() error {
	if tx.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}

	if !tx.Signature.Verify(tx.From, tx.Data) {
		return fmt.Errorf("invalid transaction signature")
	}

	return nil
}
