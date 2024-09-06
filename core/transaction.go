package core

import (
	"fmt"

	"github.com/akhilparakka/chainx/crypto"
)

type Transaction struct {
	Data []byte

	From      crypto.PublicKey
	Signature *crypto.Signature
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
