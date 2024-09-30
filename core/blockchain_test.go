package core

import (
	"testing"

	"github.com/akhilparakka/chainx/types"
	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
)

func TestAddBlock(t *testing.T) {
	bc := newblockchainWithgenesis(t)
	lenBlocks := 1000
	for i := 0; i < lenBlocks; i++ {
		block := randomblock(t, uint32(i+1), getPrevBlockhash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(block))
	}

	assert.Equal(t, bc.Height(), uint32(lenBlocks))
	assert.Equal(t, len(bc.headers), lenBlocks+1)
	assert.NotNil(t, bc.AddBlock(randomblock(t, 89, types.Hash{})))

}

func TestNewBlockchain(t *testing.T) {
	bc := newblockchainWithgenesis(t)
	assert.NotNil(t, bc.validator)
	assert.Equal(t, bc.Height(), uint32(0))
}

func TestHasBlock(t *testing.T) {
	bc := newblockchainWithgenesis(t)
	assert.True(t, bc.HasBlock(0))
	assert.False(t, bc.HasBlock(1))
	assert.False(t, bc.HasBlock(100))
}

func TestGetHeader(t *testing.T) {
	bc := newblockchainWithgenesis(t)
	lenBlocks := 1000
	for i := 0; i < lenBlocks; i++ {
		block := randomblock(t, uint32(i+1), getPrevBlockhash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(block))
		header, err := bc.GetHeader(uint32(i + 1))
		assert.Nil(t, err)
		assert.Equal(t, header, block.Header)
	}
}

func TestAddBlockToHeight(t *testing.T) {
	bc := newblockchainWithgenesis(t)

	assert.Nil(t, bc.AddBlock(randomblock(t, 1, getPrevBlockhash(t, bc, uint32(1)))))
	assert.NotNil(t, bc.AddBlock(randomblock(t, 3, types.Hash{})))
}

func newblockchainWithgenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockchain(log.NewNopLogger(), randomblock(t, 0, types.Hash{}))
	assert.Nil(t, err)
	return bc
}

func getPrevBlockhash(t *testing.T, bc *Blockchain, height uint32) types.Hash {
	prevheader, err := bc.GetHeader(height - 1)
	assert.Nil(t, err)

	return BlockHasher{}.Hash(prevheader)
}
