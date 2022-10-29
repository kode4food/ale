package types

import "fmt"

// ListOf declares a new PairType that will only allow a AnyList with elements of
// the provided elem Type
func ListOf(elem Type) Type {
	return sequenceOf(AnyList, elem)
}

// VectorOf declares a new VectorType that will only allow a AnyVector with
// elements of the provided elem Type
func VectorOf(elem Type) Type {
	return sequenceOf(AnyVector, elem)
}

func sequenceOf(base BasicType, elem Type) Type {
	name := fmt.Sprintf("%s(%s)", base.Name(), elem.Name())
	first := &namedPair{
		name: name,
		pair: &pair{
			BasicType: base,
			car:       elem,
		},
	}
	rest := &namedUnion{
		name: name,
		union: &union{
			BasicType: base,
			options:   typeList{first, Null},
		},
	}
	first.cdr = rest
	return rest
}
