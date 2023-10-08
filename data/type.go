package data

import "github.com/kode4food/ale/internal/types"

type Type struct {
	types.Type
}

func TypeOf(v Value) types.Type {
	if v, ok := v.(Typed); ok {
		return v.Type()
	}
	return types.BasicAny
}

func (t *Type) Name() Local {
	return Local(t.Type.Name())
}

func (t *Type) Call(args ...Value) Value {
	other := TypeOf(args[0])
	return Bool(types.Accepts(t.Type, other))
}

func (t *Type) Convention() Convention {
	return ApplicativeCall
}

func (t *Type) CheckArity(argCount int) error {
	return MakeFixedChecker(1)(argCount)
}

func (t *Type) Equal(other Value) bool {
	if other, ok := other.(*Type); ok {
		if t == other || t.Type == other.Type {
			return true
		}
		return t.Type.Equal(other.Type)
	}
	return false
}

func (t *Type) String() string {
	return t.Type.Name()
}
