package prototype

import (
	"crypto/sha256"
	"errors"
	"github.com/golang/protobuf/proto"
)

func (p *SignedTransaction) Id() (*Sha256, error) {
	buf, err := proto.Marshal(p)
	if err != nil {
		return nil, err
	}
	h := sha256.New()
	h.Reset()
	h.Write(buf)
	bs := h.Sum(nil)
	if bs == nil {
		return nil, errors.New("sha256 error")
	}
	id := &Sha256{Hash: bs}
	return id, nil
}
