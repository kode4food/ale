package types

type (
	Checker struct {
		parent   *Checker
		receiver Type
	}

	AcceptsWith func(*Checker, Type) bool
)

// types:
//    * any
//    * basic (string, list, cons, etc...)
//    * union
//    * sequence (vector, cons, list)
//
// accept logic:
//    * any accepts everything ✔
//    * basic left can accept itself ✔
//    * basic left can accept extended right ✔
//    * union left can accept a union subset ✔
//    * union left can accept extended union subset ✔
//    * extended left can accept itself
//    * all other permutations are false

func Accepts(left, right Type) bool {
	return Check(left).Accepts(right)
}

func accepts(c *Checker, right Type) bool {
	return c.Receiver().Accepts(c, right)
}

func Check(t Type) *Checker {
	return &Checker{
		receiver: t,
	}
}

func (c *Checker) willCycleOn(t Type) bool {
	if c.receiver == t {
		return true
	}
	if c.parent == nil {
		return false
	}
	return c.parent.willCycleOn(t)
}

func (c *Checker) Receiver() Type {
	return c.receiver
}

func (c *Checker) Accepts(right Type) bool {
	return c.AcceptsWith(right, accepts)
}

func (c *Checker) AcceptsWith(right Type, with AcceptsWith) bool {
	return with(c, right)
}

func (c *Checker) AcceptsChild(left, right Type) bool {
	return c.AcceptsChildWith(left, right, accepts)
}

func (c *Checker) AcceptsChildWith(left, right Type, with AcceptsWith) bool {
	if c.willCycleOn(left) {
		return true
	}
	return with(&Checker{
		parent:   c,
		receiver: left,
	}, right)
}
