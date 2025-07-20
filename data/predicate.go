package data

import (
	"github.com/kode4food/ale"
	"github.com/kode4food/ale/internal/types"
)

type TypePredicate struct {
	typ ale.Type
}

var typePredicateType = types.MakeBasic("type-predicate")

// MakeTypePredicate returns a TypePredicate for the given Type
func MakeTypePredicate(t ale.Type) *TypePredicate {
	return &TypePredicate{typ: t}
}

// TypePredicateOf returns a TypePredicate for the Types of the given Values.
// If more than one Value is provided, the Union of their Types is returned
func TypePredicateOf(f ale.Value, r ...ale.Value) *TypePredicate {
	if len(r) == 0 {
		return MakeTypePredicate(typeOf(f))
	}
	t := make([]ale.Type, len(r))
	for i, v := range r {
		t[i] = typeOf(v)
	}
	return MakeTypePredicate(types.MakeUnion(typeOf(f), t...))
}

func (t *TypePredicate) Type() ale.Type {
	return typePredicateType
}

func (t *TypePredicate) Name() Name {
	return Name(t.typ.Name())
}

func (t *TypePredicate) Call(args ...ale.Value) ale.Value {
	other := typeOf(args[0])
	return Bool(t.typ.Accepts(other))
}

func (t *TypePredicate) CheckArity(argc int) error {
	return CheckFixedArity(1, argc)
}

func (t *TypePredicate) Equal(other ale.Value) bool {
	if other, ok := other.(*TypePredicate); ok {
		if t == other || t.typ == other.typ {
			return true
		}
		return t.typ.Equal(other.typ)
	}
	return false
}

func (t *TypePredicate) Get(key ale.Value) (ale.Value, bool) {
	return DumpMapped(t).Get(key)
}

func typeOf(v ale.Value) ale.Type {
	if v, ok := v.(ale.Typed); ok {
		return v.Type()
	}
	return types.BasicAny
}
