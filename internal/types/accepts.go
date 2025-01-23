package types

type (
	Checker struct {
		parent   *Checker
		receiver Type
	}

	AcceptsWith func(*Checker, Type) bool
)

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
