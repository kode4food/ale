package eval

import "github.com/kode4food/ale/pkg/data"

type include struct {
	forms data.Sequence
}

func Include(forms data.Sequence) data.Value {
	return &include{forms}
}

func (i *include) Equal(other data.Value) bool {
	return i == other
}
