package prototype

import (
	"crypto/sha256"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
)

const Size = 32

func (sbh *SignedBlockHeader) Hash() (hash [Size]byte) {
	data, _ := proto.Marshal(sbh)
	hash = sha256.Sum256(data)
	return
}

func (sbh *SignedBlockHeader) Number() uint64 {
	var ret, prev BlockID
	copy(prev.Data[:], sbh.Header.Previous.Hash[:32])
	copy(ret.Data[:], sbh.BlockProducerSignature.Sig[:32])
	binary.LittleEndian.PutUint64(ret.Data[:8], prev.BlockNum()+1)
	return binary.LittleEndian.Uint64(ret.Data[:8])
}
