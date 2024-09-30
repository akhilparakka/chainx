package core

import (
	"bytes"
	"testing"
	"time"

	"github.com/akhilparakka/chainx/crypto"
	"github.com/akhilparakka/chainx/types"
	"github.com/stretchr/testify/assert"
)

func TestSignBlock(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := randomblock(t, 0, types.Hash{})
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

func TestEncodeblock(t *testing.T) {
	b := randomblock(t, 1, types.Hash{})
	buf := &bytes.Buffer{}
	assert.Nil(t, b.Encode(NewGobBlockEncoder(buf)))

	bDecode := new(Block)
	assert.Nil(t, bDecode.Decode(NewGobBlockDecoder(buf)))
	assert.Equal(t, b, bDecode)
}

func randomblock(t *testing.T, height uint32, prevblockHash types.Hash) *Block {
	privKey := crypto.GeneratePrivateKey()
	tx := randomTxWithSignature()
	header := &Header{
		Version:       1,
		PrevblockHash: prevblockHash,
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}

	b, err := Newblock(header, []*Transaction{tx})
	assert.Nil(t, err)
	dataHash, err := CalculateDataHash(b.Transactions)
	assert.Nil(t, err)
	b.Header.DataHash = dataHash
	b.Sign(privKey)

	return b
}
