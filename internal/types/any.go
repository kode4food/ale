package types

// Any accepts a Value of any other Type
type Any struct{ *Basic }

var BasicAny = &Any{
	Basic: MakeBasic("any"),
}

func (*Any) Accepts(*Checker, Type) bool {
	return true
}

func (a *Any) Equal(other Type) bool {
	_, ok := other.(*Any)
	return ok
}
