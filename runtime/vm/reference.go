package vm

import (
	"fmt"

	"github.com/kode4food/ale/data"
)

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

func (r *Ref) String() string {
	if r.Value != nil {
		return fmt.Sprintf("(ref %s)", data.MaybeQuoteString(r.Value))
	}
	return "(ref)"
}
