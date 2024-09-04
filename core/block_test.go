package core

import (
	"testing"
	"time"

	"github.com/akhilparakka/chainx/crypto"
	"github.com/akhilparakka/chainx/types"
	"github.com/stretchr/testify/assert"
)

func randomblock(height uint32) *Block {
	header := &Header{
		Version:       1,
		PrevblockHash: types.RandomHash(),
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}
	tx := Transaction{
		Data: []byte("foo"),
	}

	return Newblock(header, []Transaction{tx})
}

func Test(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := randomblock(0)
	b.Sign(privKey)

	assert.Equal(t, b.Validator, privKey.PublicKey())
	assert.NotNil(t, b.Signature)

	assert.Nil(t, b.Verify())

	attackerPrivateKey := crypto.GeneratePrivateKey()

	b.Validator = attackerPrivateKey.PublicKey()
	assert.NotNil(t, b.Verify())

	b.Sign(privKey)
	b.Height = 100
	assert.NotNil(t, b.Verify())
}
