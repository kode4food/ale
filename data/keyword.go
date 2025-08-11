package data

import (
	"fmt"
	"math/rand/v2"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/internal/types"
)

// Keyword is a Value that represents a name that resolves to itself
type Keyword string

var (
	kwdSalt = rand.Uint64()

	// compile-time checks for interface implementation
	_ interface {
		Hashed
		Procedure
		ale.Typed
		fmt.Stringer
	} = Keyword("")
)

func (k Keyword) Call(args ...ale.Value) ale.Value {
	m := args[0].(Mapped)
	res, ok := m.Get(k)
	if !ok && len(args) > 1 {
		return args[1]
	}
	return res
}

func (k Keyword) CheckArity(argc int) error {
	return CheckRangedArity(1, 2, argc)
}

func (k Keyword) Equal(other ale.Value) bool {
	return k == other
}

func (k Keyword) String() string {
	return lang.KwdPrefix + string(k)
}

func (k Keyword) Type() ale.Type {
	return types.MakeLiteral(types.BasicKeyword, k)
}

func (k Keyword) HashCode() uint64 {
	return kwdSalt ^ HashString(string(k))
}
