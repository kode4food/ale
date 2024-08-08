package data

import (
	"encoding/binary"
	"hash/maphash"
	"unsafe"
)

// Hashed can return a hash code for the value
type Hashed interface {
	HashCode() uint64
}

var seed = maphash.MakeSeed()

// HashCode returns a hash code for the provided Value. If the Value implements
// the Hashed interface, it will call us the HashCode() method. Otherwise, it
// will create a hash code from the stringified form of the Value
func HashCode(v Value) uint64 {
	if h, ok := v.(Hashed); ok {
		return h.HashCode()
	}
	return HashString(ToString(v))
}

// HashString returns a hash code for the provided string
func HashString(s string) uint64 {
	return HashBytes([]byte(s))
}

// HashInt returns a hash code for the provided int
func HashInt(i int) uint64 {
	return HashInt64(int64(i))
}

// HashInt64 returns a hash code for the provided int64
func HashInt64(i int64) uint64 {
	u := uint64(i)
	b := make([]byte, unsafe.Sizeof(u))
	binary.NativeEndian.PutUint64(b, u)
	return HashBytes(b)
}

// HashBytes returns a hash code for the provided byte slice
func HashBytes(b []byte) uint64 {
	var h maphash.Hash
	h.SetSeed(seed)
	_, _ = h.Write(b)
	return h.Sum64()
}
