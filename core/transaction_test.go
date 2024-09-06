package core

import (
	"testing"

	"github.com/akhilparakka/chainx/crypto"
	"github.com/stretchr/testify/assert"
)

func TestSignTransaction(t *testing.T) {
	privkey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: []byte("foo"),
	}

	tx.Sign(privkey)
	assert.NotNil(t, tx.Signature)
}

func TestVerifyTransaction(t *testing.T) {
	privkey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: []byte("foo"),
	}
	tx.Sign(privkey)
	assert.Nil(t, tx.Verify())

	otherprivKey := crypto.GeneratePrivateKey()
	tx.From = otherprivKey.PublicKey()

	assert.NotNil(t, tx.Verify())

}

func randomTxWithSignature() *Transaction {
	privKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: []byte("foo"),
	}
	tx.Sign(privKey)

	return tx
}
