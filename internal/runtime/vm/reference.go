package vm

import (
	"github.com/kode4food/ale"
)

// Ref encapsulates a reference to a Value
type Ref struct {
	ale.Value
}

func (r *Ref) Equal(other ale.Value) bool {
	if other, ok := other.(*Ref); ok {
		return r == other || r.Value.Equal(other.Value)
	}
	return r.Value.Equal(other)
}
