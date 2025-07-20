package sequence

import (
	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
)

func Take(s data.Sequence, count int) (data.Vector, data.Sequence, bool) {
	var f ale.Value
	var ok bool
	res := make(data.Vector, count)
	for i := range count {
		if f, s, ok = s.Split(); !ok {
			return data.EmptyVector, data.Null, false
		}
		res[i] = f
	}
	return res, s, true
}
