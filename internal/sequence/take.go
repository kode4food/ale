package sequence

import "github.com/kode4food/ale/pkg/data"

func Take(s data.Sequence, count int) (data.Vector, data.Sequence, bool) {
	var f data.Value
	var ok bool
	res := make(data.Vector, count)
	for i := 0; i < count; i++ {
		if f, s, ok = s.Split(); !ok {
			return data.EmptyVector, data.Null, false
		}
		res[i] = f
	}
	return res, s, true
}
