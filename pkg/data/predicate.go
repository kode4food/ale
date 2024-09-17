package data

import "github.com/kode4food/ale/internal/types"

type TypePredicate struct {
	typ types.Type
}

var typePredicateType = types.MakeBasic("type-predicate")

// MakeTypePredicate returns a TypePredicate for the given Type
func MakeTypePredicate(t types.Type) *TypePredicate {
	return &TypePredicate{typ: t}
}

// TypePredicateOf returns a TypePredicate for the Types of the given Values.
// If more than one Value is provided, the Union of their Types will be
// returned
func TypePredicateOf(f Value, r ...Value) *TypePredicate {
	if len(r) == 0 {
		return MakeTypePredicate(typeOf(f))
	}
	t := make([]types.Type, len(r))
	for i, v := range r {
		t[i] = typeOf(v)
	}
	return MakeTypePredicate(types.MakeUnion(typeOf(f), t...))
}

func (t *TypePredicate) Type() types.Type {
	return typePredicateType
}

func (t *TypePredicate) Name() Local {
	return Local(t.typ.Name())
}

func (t *TypePredicate) Call(args ...Value) Value {
	other := typeOf(args[0])
	return Bool(types.Accepts(t.typ, other))
}

func (t *TypePredicate) CheckArity(argc int) error {
	return checkFixedArity(1, argc)
}

func (t *TypePredicate) Equal(other Value) bool {
	if other, ok := other.(*TypePredicate); ok {
		if t == other || t.typ == other.typ {
			return true
		}
		return t.typ.Equal(other.typ)
	}
	return false
}

func (t *TypePredicate) Get(key Value) (Value, bool) {
	return DumpMapped(t).Get(key)
}

func typeOf(v Value) types.Type {
	if v, ok := v.(Typed); ok {
		return v.Type()
	}
	return types.BasicAny
}
