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
	if a == other {
		return true
	}
	if _, ok := other.(*Any); ok {
		return true
	}
	return false
}
