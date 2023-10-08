package types

import "fmt"

// MakeTuple declares a new TupleType that will only allow a MakeListOf or MakeVectorOf
// with positional elements of the provided Types
func MakeTuple(elems ...Type) Type {
	base := []*Basic{BasicCons, BasicList, BasicVector}
	res := make([]Type, 3)
	for idx, t := range base {
		var comp Type = BasicNull
		for i := len(elems) - 1; i >= 0; i = i - 1 {
			comp = &Pair{
				basic: t,
				name:  fmt.Sprintf("tuple(%s)", typeList(elems).name()),
				car:   elems[i],
				cdr:   comp,
			}
		}
		res[idx] = comp
	}
	return &Union{
		name:    fmt.Sprintf("tuple(%s)", typeList(elems).name()),
		basic:   BasicUnion,
		options: res,
	}
}
