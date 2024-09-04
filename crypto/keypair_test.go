package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateKeys(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.PublicKey()

	msg := []byte("Hello World")

	sig := privKey.Sign(msg)
	assert.True(t, sig.Verify(pubKey, msg))
}

func TestKeypairSignVerifyFail(t *testing.T) {
	privkey := GeneratePrivateKey()
	msg := []byte("Hello World")

	sig := privkey.Sign(msg)

	otherPrivateKey := GeneratePrivateKey()
	otherpubKey := otherPrivateKey.PublicKey()

	assert.False(t, sig.Verify(otherpubKey, msg))
	assert.False(t, sig.Verify(privkey.PublicKey(), []byte("Hello World!")))
}
