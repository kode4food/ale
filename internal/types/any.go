package types

import "github.com/kode4food/ale"

// Any accepts a Value of any other Type
type Any struct{ *Basic }

var BasicAny = &Any{
	Basic: makeBasic("any"),
}

func (*Any) Accepts(ale.Type) bool {
	return true
}

func (a *Any) Equal(other ale.Type) bool {
	_, ok := other.(*Any)
	return ok
}
