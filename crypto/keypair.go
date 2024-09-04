package crypto

import (
	"crypto/sha256"

	"github.com/akhilparakka/chainx/types"
	"github.com/cloudflare/circl/sign/dilithium"
)

type PrivateKey struct {
	key dilithium.PrivateKey
}

func (k PrivateKey) Sign(data []byte) *Signature {
	mode := dilithium.Mode3
	sign := mode.Sign(k.key, data)
	return &Signature{
		signature: sign,
	}
}

type PublicKey struct {
	key dilithium.PublicKey
}

func GeneratePrivateKey() PrivateKey {
	mode := dilithium.Mode3

	_, privKey, err := mode.GenerateKey(nil)
	if err != nil {
		panic(err)
	}

	return PrivateKey{
		key: privKey,
	}
}

func (k PrivateKey) PublicKey() PublicKey {
	return PublicKey{
		key: k.key.Public().(dilithium.PublicKey),
	}
}

func (k PublicKey) ToSlice() []byte {
	return k.key.Bytes()
}

func (k PublicKey) Address() types.Address {
	h := sha256.Sum256(k.ToSlice())
	return types.AddressFromBytes(h[len(h)-20:])
}

type Signature struct {
	signature []byte
}

func (sig Signature) Verify(publicKey PublicKey, data []byte) bool {
	mode := dilithium.Mode3
	return mode.Verify(publicKey.key, data, sig.signature)
}
