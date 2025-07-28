package types

import (
	"fmt"
	"slices"

	"github.com/kode4food/ale"
)

// Union describes a Type that can accept any of a set of Types
type Union struct {
	basic
	name    string
	options typeList
}

// MakeUnion declares a *Union based on at least one provided Type. If any of
// the provided types is a types.Any, then types.Any will be returned
func MakeUnion(first ale.Type, rest ...ale.Type) ale.Type {
	if len(rest) == 0 {
		return first
	}
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

func (u *Union) Options() []ale.Type {
	return u.options
}

func (u *Union) Accepts(other ale.Type) bool {
	switch other := other.(type) {
	case *Union:
		return u == other || compoundAccepts(u, other)
	default:
		return compoundAccepts(u, other)
	}
}

func (u *Union) accepts(c *checker, other ale.Type) bool {
	switch other := other.(type) {
	case *Union:
		return u == other || u.acceptsUnion(c, other)
	default:
		return u.acceptsType(c, other)
	}
}

func (u *Union) Equal(other ale.Type) bool {
	if other, ok := other.(*Union); ok {
		return u == other ||
			u.name == other.name &&
				u.basic.Equal(other.basic) &&
				u.options.equal(other.options)
	}
	return false
}

func (u *Union) acceptsUnion(c *checker, other *Union) bool {
	for _, o := range other.Options() {
		if !u.acceptsType(c, o) {
			return false
		}
	}
	return true
}

func (u *Union) acceptsType(c *checker, other ale.Type) bool {
	return slices.ContainsFunc(u.options, func(t ale.Type) bool {
		return c.acceptsChild(t, other)
	})
}
