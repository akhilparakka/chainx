package core

import (
	"encoding/gob"
	"io"

	"github.com/cloudflare/circl/sign/dilithium"
	"github.com/cloudflare/circl/sign/dilithium/mode3"
)

type Encoder[T any] interface {
	Encode(T) error
}

type Decoder[T any] interface {
	Decode(T) error
}

type GobTxEncoder struct {
	w io.Writer
}

func NewGobTxEncoder(w io.Writer) *GobTxEncoder {
	return &GobTxEncoder{
		w: w,
	}
}

func (e *GobTxEncoder) Encode(tx *Transaction) error {
	gob.Register((*dilithium.PublicKey)(nil))
	gob.Register((*mode3.PublicKey)(nil))
	return gob.NewEncoder(e.w).Encode(tx)
}

type GobTxDecoder struct {
	r io.Reader
}

func NewGobTxDecoder(r io.Reader) *GobTxDecoder {
	return &GobTxDecoder{
		r: r,
	}
}

func (e *GobTxDecoder) Decode(tx *Transaction) error {
	gob.Register((*dilithium.PublicKey)(nil))
	gob.Register((*mode3.PublicKey)(nil))
	gob.Register(tx)
	return gob.NewDecoder(e.r).Decode(tx)
}

type GobBlockEncoder struct {
	w io.Writer
}

func NewGobBlockEncoder(w io.Writer) *GobBlockEncoder {
	return &GobBlockEncoder{
		w: w,
	}
}

func (enc *GobBlockEncoder) Encode(b *Block) error {
	return gob.NewEncoder(enc.w).Encode(b)
}

type GobBlockDecoder struct {
	r io.Reader
}

func NewGobBlockDecoder(r io.Reader) *GobBlockDecoder {
	return &GobBlockDecoder{
		r: r,
	}
}

func (enc *GobBlockDecoder) Decode(b *Block) error {
	return gob.NewDecoder(enc.r).Decode(b)
}
