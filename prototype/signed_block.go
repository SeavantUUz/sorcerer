package prototype

import "encoding/binary"

type BlockID struct {
	Data [32]byte
}

func (bid BlockID) BlockNum() uint64 {
	return binary.LittleEndian.Uint64(bid.Data[:8])
}

func (sb *SignedBlock) Id() BlockID {
	var ret, prev BlockID
	if sb.SignedHeader != nil && sb.SignedHeader.Header != nil &&
		sb.SignedHeader.Header.Previous != nil &&
		len(sb.SignedHeader.Header.Previous.Hash) != 0 {
		copy(prev.Data[:], sb.SignedHeader.Header.Previous.Hash[:32])
	}
	digest := sb.SignedHeader.Hash()
	copy(ret.Data[:], digest[:])
	binary.LittleEndian.PutUint64(ret.Data[:8], prev.BlockNum()+1)
	return ret
}
