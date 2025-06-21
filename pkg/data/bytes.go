package data

import (
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"strings"

	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/internal/types"
)

type Bytes []byte

const (
	ErrIntegerOutOfRange = "integer out of byte range: %d"
)

var (
	EmptyBytes Bytes

	bytesSalt = rand.Uint64()
)

func NewBytes(vals ...Value) Bytes {
	res, err := ValuesToBytes(vals...)
	if err != nil {
		panic(err)
	}
	return res
}

func ValuesToBytes(vals ...Value) (Bytes, error) {
	if len(vals) == 0 {
		return EmptyBytes, nil
	}
	res := make(Bytes, len(vals))
	for i, v := range vals {
		b, err := toByte(v)
		if err != nil {
			return nil, err
		}
		res[i] = byte(b)
	}
	return res, nil
}

func (b Bytes) Count() int {
	return len(b)
}

func (b Bytes) ElementAt(index int) (Value, bool) {
	if index >= 0 && index < len(b) {
		return Integer(b[index]), true
	}
	return Null, false
}

func (b Bytes) IsEmpty() bool {
	return len(b) == 0
}

func (b Bytes) Car() Value {
	if len(b) > 0 {
		return Integer(b[0])
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
		return Integer(b[0]), EmptyBytes, true
	default:
		return Integer(b[0]), b[1:], true
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
	return Bytes(sliceRangedCall(b, args))
}

func (b Bytes) CheckArity(argc int) error {
	return CheckRangedArity(1, 2, argc)
}

func (b Bytes) Equal(other Value) bool {
	if o, ok := other.(Bytes); ok {
		return slices.Equal(b, o)
	}
	return false
}

func (b Bytes) String() string {
	var buf strings.Builder
	buf.WriteString(lang.BytesStart)
	for i, v := range b {
		if i > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(Integer(v).String())
	}
	buf.WriteString(lang.BytesEnd)
	return buf.String()
}

func (Bytes) Type() types.Type {
	return types.BasicBytes
}

func (b Bytes) HashCode() uint64 {
	return bytesSalt ^ HashBytes(b)
}

func mustToByte(v Value) byte {
	b, err := toByte(v)
	if err != nil {
		panic(err)
	}
	return b
}

func toByte(v Value) (byte, error) {
	if v, ok := v.(Integer); ok {
		if v < 0 || v > math.MaxUint8 {
			return 0, fmt.Errorf(ErrIntegerOutOfRange, v)
		}
		return byte(v), nil
	}
	return 0, fmt.Errorf(ErrExpectedInteger, v)
}
