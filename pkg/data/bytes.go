package data

import (
	"fmt"
	"math"
	"math/rand/v2"
	"slices"

	"github.com/kode4food/ale/internal/types"
)

type (
	Byte  byte
	Bytes []byte
)

var (
	EmptyBytes Bytes

	bytesSalt = rand.Uint64()
	bytesType = types.MakeBasic("bytes")
)

func (b Byte) Equal(other Value) bool {
	if o, ok := other.(Byte); ok {
		return b == o
	}
	return false
}

func (b Byte) String() string {
	return fmt.Sprintf("%d", b)
}

func (b Bytes) Count() int {
	return len(b)
}

func (b Bytes) ElementAt(index int) (Value, bool) {
	if index >= 0 && index < len(b) {
		return Byte(b[index]), true
	}
	return Null, false
}

func (b Bytes) IsEmpty() bool {
	return len(b) == 0
}

func (b Bytes) Car() Value {
	if len(b) > 0 {
		return Byte(b[0])
	}
	return Null
}

func (b Bytes) Cdr() Value {
	if len(b) > 1 {
		return b[1:]
	}
	return EmptyBytes
}

func (b Bytes) Split() (Value, Sequence, bool) {
	switch len(b) {
	case 0:
		return Null, EmptyBytes, false
	case 1:
		return Byte(b[0]), EmptyBytes, true
	default:
		return Byte(b[0]), b[1:], true
	}
}

func (b Bytes) Prepend(v Value) Sequence {
	vb := mustToByte(v)
	res := make(Bytes, len(b)+1)
	res[0] = byte(vb)
	copy(res[1:], b)
	return res
}

func (b Bytes) Append(v Value) Sequence {
	vb := mustToByte(v)
	res := make(Bytes, len(b)+1)
	copy(res, b)
	res[len(b)] = byte(vb)
	return res
}

func (b Bytes) Reverse() Sequence {
	vl := len(b)
	if vl <= 1 {
		return b
	}
	res := make(Bytes, vl)
	for i, j := 0, vl-1; j >= 0; i, j = i+1, j-1 {
		res[i] = b[j]
	}
	return res
}

func (b Bytes) IndexOf(v Value) (int, bool) {
	vb := mustToByte(v)
	i := slices.Index(b, byte(vb))
	return i, i != -1
}

func (b Bytes) Call(args ...Value) Value {
	if len(args) == 1 {
		return b.callFrom(int(args[0].(Integer)))
	}
	return b.callRange(int(args[0].(Integer)), int(args[1].(Integer)))
}

func (b Bytes) callFrom(idx int) Value {
	if idx < 0 {
		panic(fmt.Errorf(ErrInvalidStartIndex, idx))
	}
	if ns, ok := b.from(idx); ok {
		return ns
	}
	panic(fmt.Errorf(ErrInvalidStartIndex, idx))
}

func (b Bytes) callRange(idx, end int) Value {
	if idx < 0 || end < idx {
		panic(fmt.Errorf(ErrInvalidIndexes, idx, end))
	}

	ns, ok := b.from(idx)
	if !ok || len(ns) == 0 && end > idx {
		panic(fmt.Errorf(ErrInvalidIndexes, idx, end))
	}

	if res, ok := ns.take(end - idx); ok {
		return res
	}
	panic(fmt.Errorf(ErrInvalidIndexes, idx, end))
}

func (b Bytes) CheckArity(argc int) error {
	return CheckRangedArity(1, 2, argc)
}

func (b Bytes) from(idx int) (Bytes, bool) {
	_, r, ok := b.splitAt(idx)
	return r, ok
}

func (b Bytes) take(count int) (Bytes, bool) {
	f, _, ok := b.splitAt(count)
	return f, ok
}

func (b Bytes) splitAt(idx int) (Bytes, Bytes, bool) {
	if idx >= 0 && idx < len(b) {
		return b[:idx], b[idx:], true
	}
	return EmptyBytes, EmptyBytes, false
}

func (b Bytes) Equal(other Value) bool {
	if o, ok := other.(Bytes); ok {
		return slices.Equal(b, o)
	}
	return false
}

func (Bytes) Type() types.Type {
	return bytesType
}

func (b Bytes) HashCode() uint64 {
	return bytesSalt ^ HashBytes(b)
}

func mustToByte(v Value) Byte {
	if b, ok := toByte(v); ok {
		return b
	}
	panic(fmt.Errorf("not a valid byte value: %s", v))
}

func toByte(v Value) (Byte, bool) {
	switch e := v.(type) {
	case Byte:
		return e, true
	case Integer:
		if e < 0 || e > math.MaxUint8 {
			panic("integer out of byte range")
		}
		return Byte(e), true
	default:
		return 0, false
	}
}
