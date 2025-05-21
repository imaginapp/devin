package hash

import (
	"bytes"
	"encoding/binary"
	"io"
	"strconv"

	"github.com/OneOfOne/xxhash"
)

// HashSeed -
var HashSeed uint64 = 736626791

// Hash the hashed response
type Hash struct {
	Sum64 uint64
}

// Base36 return hash in base36 string
func (h *Hash) Base36() string {
	return strconv.FormatUint(h.Sum64, 36)
}

func (h *Hash) String() string {
	return h.Base36()
}

// Bytes returns hash in bytes
func (h *Hash) Bytes() []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(h.Sum64))
	return b
}

// Reader returns hash in a io.Reader
func (h *Hash) Reader() io.Reader {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], h.Sum64)
	return bytes.NewReader(buf[:])
}

// StringHash -
func StringHash(s string) *Hash {
	sum64 := xxhash.ChecksumString64S(s, HashSeed)
	return &Hash{Sum64: sum64}
}

// BytesHash -
func BytesHash(b []byte) *Hash {
	sum64 := xxhash.Checksum64S(b, HashSeed)
	return &Hash{Sum64: sum64}
}
