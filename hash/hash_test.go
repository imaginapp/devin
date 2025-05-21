package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	HashSeed = 1234
	assert := assert.New(t)
	hash := StringHash("zdpuB1gdDa4uHqqkYC2AJ9Nd8yhrnLZupC83SzDDaMxG4JzjZ")
	expected := "27595fkgr1xkm"
	assert.Equal(expected, hash.String())
}

func TestHashBytesSum(t *testing.T) {
	HashSeed = 1234
	assert := assert.New(t)
	hash := BytesHash([]byte("Some data"))
	expected := uint64(11072410196304768542)
	assert.Equal(expected, hash.Sum64)
}
