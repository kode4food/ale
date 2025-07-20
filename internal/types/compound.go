package types

import "github.com/kode4food/ale"

type (
	cycleChecker struct {
		parent   *cycleChecker
		receiver ale.Type
	}

	// compound is an internal interface for types that need cycle checking
	compound interface {
		accepts(*cycleChecker, ale.Type) bool
	}
)

func compoundAccepts(left, right ale.Type) bool {
	c := &cycleChecker{
		receiver: left,
	}
	return c.accepts(right)
}

func (c *cycleChecker) accepts(right ale.Type) bool {
	if r, ok := c.receiver.(compound); ok {
		return r.accepts(c, right)
	}
	return c.receiver.Accepts(right)
}

func (c *cycleChecker) willCycleOn(t ale.Type) bool {
	if c.receiver == t {
		return true
	}
	if c.parent == nil {
		return false
	}
	return c.parent.willCycleOn(t)
}

func (c *cycleChecker) acceptsChild(left, right ale.Type) bool {
	if c.willCycleOn(left) {
		return true
	}
	child := &cycleChecker{
		parent:   c,
		receiver: left,
	}
	return child.accepts(right)
}
