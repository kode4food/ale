package types

import (
	"fmt"

	"github.com/kode4food/ale"
)

type Literal struct {
	Basic
	value ale.Value
}

func MakeLiteral(b Basic, v ale.Value) *Literal {
	return &Literal{
		Basic: b,
		value: v,
	}
}

func (l *Literal) Type() ale.Type {
	return l.Basic
}

func (l *Literal) Value() ale.Value {
	return l.value
}

func (l *Literal) Name() string {
	if v, ok := l.value.(fmt.Stringer); ok {
		return fmt.Sprintf("%s(%s)", l.Basic.Name(), v)
	}
	return fmt.Sprintf("%s(%p)", l.Basic.Name(), l.value)
}

func (l *Literal) Accepts(other ale.Type) bool {
	if other, ok := other.(*Literal); ok {
		return l.Basic.Accepts(other.Basic) && l.value.Equal(other.value)
	}
	return false
}

func (l *Literal) Equal(other ale.Type) bool {
	if other, ok := other.(*Literal); ok {
		return l.Basic.Equal(other.Basic) && l.value.Equal(other.value)
	}
	return false
}
