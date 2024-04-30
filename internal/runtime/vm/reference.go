package vm

import "github.com/kode4food/ale/pkg/data"

// Ref encapsulates a reference to a Value
type Ref struct {
	data.Value
}

func (r *Ref) Equal(other data.Value) bool {
	if other, ok := other.(*Ref); ok {
		return r.Value.Equal(other.Value)
	}
	return r.Value.Equal(other)
}
