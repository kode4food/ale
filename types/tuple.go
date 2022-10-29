package types

import "fmt"

// Tuple declares a new TupleType that will only allow a ListOf or VectorOf
// with positional elements of the provided Types
func Tuple(elems ...Type) Type {
	base := []BasicType{AnyCons, AnyList, AnyVector}
	res := make([]Type, 3)
	for idx, t := range base {
		var comp Type = Null
		for i := len(elems) - 1; i >= 0; i = i - 1 {
			comp = &namedPair{
				name: fmt.Sprintf("tuple(%s)", typeList(elems).name()),
				pair: &pair{
					BasicType: t,
					car:       elems[i],
					cdr:       comp,
				},
			}
		}
		res[idx] = comp
	}
	return &namedUnion{
		name: fmt.Sprintf("tuple(%s)", typeList(elems).name()),
		union: &union{
			BasicType: AnyUnion,
			options:   res,
		},
	}
}
