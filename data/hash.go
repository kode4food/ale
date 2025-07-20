package data

import (
	"encoding/binary"
	"hash/maphash"

	"github.com/kode4food/ale"
)

// Hashed can return a hash code for the value
type Hashed interface {
	ale.Value
	HashCode() uint64
}

const int64CacheSize = 1024

var (
	seed        = maphash.MakeSeed()
	int64Hashes = makeCachedInt64Hashes()
)

// HashCode returns a hash code for the provided Value. If the Value implements
// the Hashed interface, it will call us the HashCode() method. Otherwise, it
// will create a hash code from the stringified form of the Value
func HashCode(v ale.Value) uint64 {
	if h, ok := v.(Hashed); ok {
		return h.HashCode()
	}
	return HashString(ToString(v))
}

// HashString returns a hash code for the provided string
func HashString(s string) uint64 {
	return maphash.String(seed, s)
}

// HashInt returns a hash code for the provided int
func HashInt(i int) uint64 {
	return HashInt64(int64(i))
}

// HashInt64 returns a hash code for the provided int64
func HashInt64(i int64) uint64 {
	if i < int64CacheSize {
		return int64Hashes[i]
	}
	return hashInt64(i)
}

func hashInt64(i int64) uint64 {
	var b [8]byte
	s := b[:]
	binary.NativeEndian.PutUint64(s, uint64(i))
	return maphash.Bytes(seed, s)
}

// HashBytes returns a hash code for the provided byte slice
func HashBytes(b []byte) uint64 {
	return maphash.Bytes(seed, b)
}

func makeCachedInt64Hashes() []uint64 {
	res := make([]uint64, int64CacheSize)
	for i := range res {
		res[i] = hashInt64(int64(i))
	}
	return res
}
