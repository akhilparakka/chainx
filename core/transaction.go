package core

import (
	"fmt"

	"github.com/akhilparakka/chainx/crypto"
)

type Transaction struct {
	Data []byte

	PublicKey crypto.PublicKey
	Signature *crypto.Signature
}

func (tx *Transaction) Sign(privkey crypto.PrivateKey) {
	sig := privkey.Sign(tx.Data)

	tx.PublicKey = privkey.PublicKey()
	tx.Signature = sig
}

func (tx *Transaction) Verify() error {
	if tx.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}

	if !tx.Signature.Verify(tx.PublicKey, tx.Data) {
		return fmt.Errorf("invalid transaction signature")
	}

	return nil
}
