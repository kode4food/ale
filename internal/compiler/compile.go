package compiler

import (
	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
)

// IsEvaluable returns whether the provided value is subject to further
// evaluation
func IsEvaluable(v ale.Value) bool {
	switch v.(type) {
	case data.Symbol, *data.List, data.Vector, *data.Object:
		return true
	default:
		return false
	}
}
