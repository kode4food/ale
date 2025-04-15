package types

import "fmt"

// MakeListOf declares a new PairType that will only allow a BasicList with
// elements of the provided elem Type
func MakeListOf(elem Type) Type {
	return sequenceOf(BasicList, elem)
}

// MakeVectorOf declares a new VectorType that will only allow a BasicVector
// with elements of the provided elem Type
func MakeVectorOf(elem Type) Type {
	return sequenceOf(BasicVector, elem)
}

func sequenceOf(base *Basic, elem Type) Type {
	name := fmt.Sprintf("%s(%s)", base.Name(), elem.Name())
	first := &Pair{
		basic: base,
		name:  name,
		car:   elem,
	}
	rest := &Union{
		name:    name,
		basic:   base,
		options: typeList{first, BasicNull},
	}
	first.cdr = rest
	return rest
}
