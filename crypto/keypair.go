package crypto

import (
	"crypto/sha256"

	"github.com/akhilparakka/chainx/types"
	"github.com/cloudflare/circl/sign/dilithium"
)

var mode = dilithium.Mode3

type PrivateKey struct {
	Key dilithium.PrivateKey
}

func (k PrivateKey) Sign(data []byte) *Signature {
	sign := mode.Sign(k.Key, data)
	return &Signature{
		Signature: sign,
	}
}

type PublicKey struct {
	Key dilithium.PublicKey
}

func GeneratePrivateKey() PrivateKey {
	_, privKey, err := mode.GenerateKey(nil)
	if err != nil {
		panic(err)
	}

	return PrivateKey{
		Key: privKey,
	}
}

func (k PrivateKey) PublicKey() PublicKey {
	return PublicKey{
		Key: k.Key.Public().(dilithium.PublicKey),
	}
}

func (k PublicKey) ToSlice() []byte {
	return k.Key.Bytes()
}

func (k PublicKey) Address() types.Address {
	h := sha256.Sum256(k.ToSlice())
	return types.AddressFromBytes(h[len(h)-20:])
}

type Signature struct {
	Signature []byte
}

func (sig Signature) Verify(publicKey PublicKey, data []byte) bool {
	return mode.Verify(publicKey.Key, data, sig.Signature)
}
