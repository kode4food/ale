package data

import "github.com/kode4food/ale/internal/types"

type Type struct {
	typ types.Type
}

var typePredicateType = types.MakeBasic("type-predicate")

// MakeType returns a type Predicate for the given Type
func MakeType(t types.Type) *Type {
	return &Type{typ: t}
}

// TypeOf returns a type Predicate for matching the type of the given Value
func TypeOf(v Value) *Type {
	return MakeType(typeOf(v))
}

func (t *Type) Type() types.Type {
	return typePredicateType
}

func (t *Type) Name() Local {
	return Local(t.typ.Name())
}

func (t *Type) Call(args ...Value) Value {
	other := typeOf(args[0])
	return Bool(types.Accepts(t.typ, other))
}

func (t *Type) Convention() Convention {
	return ApplicativeCall
}

func (t *Type) CheckArity(argCount int) error {
	return MakeFixedChecker(1)(argCount)
}

func (t *Type) Equal(other Value) bool {
	if other, ok := other.(*Type); ok {
		if t == other || t.typ == other.typ {
			return true
		}
		return t.typ.Equal(other.typ)
	}
	return false
}

func (t *Type) String() string {
	return DumpString(t)
}

func typeOf(v Value) types.Type {
	if v, ok := v.(Typed); ok {
		return v.Type()
	}
	return types.BasicAny
}
