package core

import (
	"testing"
	"time"

	"github.com/akhilparakka/chainx/crypto"
	"github.com/akhilparakka/chainx/types"
	"github.com/stretchr/testify/assert"
)

func TestSignBlock(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := randomblock(0, types.Hash{})
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

func randomblock(height uint32, prevblockHash types.Hash) *Block {
	header := &Header{
		Version:       1,
		PrevblockHash: prevblockHash,
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}

	return Newblock(header, []Transaction{})
}

func randomBlockWithSignature(height uint32, prevBlockHash types.Hash) *Block {
	privkey := crypto.GeneratePrivateKey()
	b := randomblock(height, prevBlockHash)
	tx := randomTxWithSignature()
	b.AddTransaction(tx)
	b.Sign(privkey)

	return b
}
