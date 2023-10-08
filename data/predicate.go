package data

import "github.com/kode4food/ale/internal/types"

type TypePredicate struct {
	typ types.Type
}

var typePredicateType = types.MakeBasic("type-predicate")

// MakeType returns a type Predicate for the given Type
func MakeType(t types.Type) *TypePredicate {
	return &TypePredicate{typ: t}
}

// TypeOf returns a TypePredicate for matching the type of the given Value
func TypeOf(v Value) *TypePredicate {
	return MakeType(typeOf(v))
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

func (t *TypePredicate) Convention() Convention {
	return ApplicativeCall
}

func (t *TypePredicate) CheckArity(argCount int) error {
	return MakeFixedChecker(1)(argCount)
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

func (t *TypePredicate) String() string {
	return DumpString(t)
}

func typeOf(v Value) types.Type {
	if v, ok := v.(Typed); ok {
		return v.Type()
	}
	return types.BasicAny
}
