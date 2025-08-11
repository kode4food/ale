package data

import (
	"math/rand/v2"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/internal/types"
)

// Bool represents the data True or False
type Bool bool

const (
	// True represents the boolean value of True
	True Bool = true

	// False represents the boolean value of false
	False Bool = false
)

var (
	trueHash  = rand.Uint64()
	falseHash = rand.Uint64()

	// compile-time checks for interface implementation
	_ interface {
		Hashed
		Procedure
		ale.Typed
	} = False
)

func (b Bool) Call(...ale.Value) ale.Value {
	return b
}

func (b Bool) CheckArity(int) error {
	return nil
}

func (b Bool) Equal(other ale.Value) bool {
	return b == other
}

func (b Bool) String() string {
	if b {
		return lang.TrueLiteral
	}
	return lang.FalseLiteral
}

func (b Bool) Type() ale.Type {
	return types.MakeLiteral(types.BasicBoolean, b)
}

func (b Bool) HashCode() uint64 {
	if b {
		return trueHash
	}
	return falseHash
}
