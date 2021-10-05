package extended

import "github.com/kode4food/ale/types"

type extended struct {
	types.Type
}

// New creates an Extended base for the specified Type
func New(t types.Type) types.Extended {
	return &extended{
		Type: t,
	}
}

func (e *extended) Base() types.Type {
	return e.Type
}
