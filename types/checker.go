package types

type (
	Checker interface {
		Check(Type) Checker
		Accepts(Type) Accepted
	}

	Accepted interface {
		Checker
	}

	checker struct {
		parent   *checker
		receiver Type
	}

	accepted struct {
		*checker
		other Type
	}
)

func Check(receiver Type) Checker {
	return &checker{
		receiver: receiver,
	}
}

func (c *checker) Check(receiver Type) Checker {
	return &checker{
		parent:   c,
		receiver: receiver,
	}
}

func (c *checker) Accepts(other Type) Accepted {
	if c.willAccept(other) {
		return &accepted{
			checker: c,
			other:   other,
		}
	}
	return nil
}

func (c *checker) willAccept(other Type) bool {
	if c.checkReceiverCycle() {
		return true
	}
	return c.receiver.Accepts(c, other)
}

func (c *checker) checkReceiverCycle() bool {
	r := c.receiver
	for a := c.parent; a != nil; a = a.parent {
		if a.receiver == r {
			return true
		}
	}
	return false
}
