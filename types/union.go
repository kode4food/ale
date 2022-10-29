package types

import "fmt"

type (
	// UnionType describes a Type that can accept any of a set of Types
	UnionType interface {
		Type
		union() // marker
		Options() []Type
	}

	union struct {
		BasicType
		options typeList
	}

	namedUnion struct {
		*union
		name string
	}
)

// Union declares a UnionType based on at least one provided Type. If any of
// the provided types is a types.Any, then types.Any will be returned
func Union(first Type, rest ...Type) Type {
	all := append(typeList{first}, rest...).flatten()
	if all.hasAny() {
		return Any
	}
	return &union{
		BasicType: unionBasicType(all),
		options:   all,
	}
}

func unionBasicType(t typeList) BasicType {
	if res, ok := t.basicType(); ok {
		return res
	}
	return AnyUnion
}

func (*union) union() {}

func (u *union) Name() string {
	return fmt.Sprintf("union(%s)", u.options.sorted().name())
}

func (u *union) Options() []Type {
	return u.options
}

func (u *union) Accepts(c *Checker, other Type) bool {
	switch other := other.(type) {
	case UnionType:
		return u.acceptsUnion(c, other)
	default:
		return u.acceptsType(c, other)
	}
}

func (u *union) acceptsUnion(c *Checker, other UnionType) bool {
	for _, o := range other.Options() {
		if !u.acceptsType(c, o) {
			return false
		}
	}
	return true
}

func (u *union) acceptsType(c *Checker, other Type) bool {
	for _, t := range u.options {
		if c.AcceptsChild(t, other) {
			return true
		}
	}
	return false
}

func (u *namedUnion) Name() string {
	return u.name
}

func (u *namedUnion) Accepts(c *Checker, other Type) bool {
	if other, ok := other.(*namedUnion); ok {
		return u.union.Accepts(c, other.union)
	}
	return u.union.Accepts(c, other)
}
