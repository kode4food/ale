package types

import "github.com/kode4food/ale"

// Any accepts a Value of any other Type
type Any struct{ *basic }

var BasicAny = &Any{
	basic: makeBasic("any"),
}

func (*Any) Accepts(ale.Type) bool {
	return true
}

func (a *Any) Equal(other ale.Type) bool {
	_, ok := other.(*Any)
	return ok
}
