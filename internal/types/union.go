package types

import (
	"fmt"
	"slices"
)

// Union describes a Type that can accept any of a set of Types
type Union struct {
	basic
	name    string
	options typeList
}

// MakeUnion declares a *Union based on at least one provided Type. If any of
// the provided types is a types.Any, then types.Any will be returned
func MakeUnion(first Type, rest ...Type) Type {
	all := append(typeList{first}, rest...).flatten()
	if all.hasAny() {
		return BasicAny
	}
	return &Union{
		basic:   unionBasicType(all),
		name:    fmt.Sprintf("union(%s)", all.sorted().name()),
		options: all,
	}
}

func unionBasicType(t typeList) basic {
	if res, ok := t.basicType(); ok {
		return res
	}
	return BasicUnion
}

func (u *Union) Name() string {
	return u.name
}

func (u *Union) Options() []Type {
	return u.options
}

func (u *Union) Accepts(c *Checker, other Type) bool {
	switch other := other.(type) {
	case *Union:
		return u == other || u.acceptsUnion(c, other)
	default:
		return u.acceptsType(c, other)
	}
}

func (u *Union) Equal(other Type) bool {
	if other, ok := other.(*Union); ok {
		return u == other ||
			u.name == other.name &&
				u.basic.Equal(other.basic) &&
				u.options.equal(other.options)
	}
	return false
}

func (u *Union) acceptsUnion(c *Checker, other *Union) bool {
	for _, o := range other.Options() {
		if !u.acceptsType(c, o) {
			return false
		}
	}
	return true
}

func (u *Union) acceptsType(c *Checker, other Type) bool {
	return slices.ContainsFunc(u.options, func(t Type) bool {
		return c.AcceptsChild(t, other)
	})
}
