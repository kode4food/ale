package builtin

import "github.com/kode4food/ale/pkg/data"

func isNaN(v data.Value) bool {
	if num, ok := v.(data.Number); ok {
		return num.IsNaN()
	}
	return false
}
