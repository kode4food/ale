package compound

import (
	"fmt"

	"github.com/kode4food/ale/types"
	"github.com/kode4food/ale/types/basic"
	"github.com/kode4food/ale/types/extended"
)

type (
	// UnionType describes a Type that can accept any of a set of Types
	UnionType interface {
		types.Extended
		union() // marker
		Options() []types.Type
	}

	union struct {
		types.Extended
		options typeList
	}
)

// Union declares a UnionType based on at least one provided Type. If any of
// the provided types is a types.Any, then types.Any will be returned
func Union(first types.Type, rest ...types.Type) types.Type {
	all := append(typeList{first}, rest...).flatten()
	if all.hasAny() {
		return basic.Any
	}
	return &union{
		Extended: extended.New(basicType(all)),
		options:  all,
	}
}

func (*union) union() {}

func (u *union) Name() string {
	return fmt.Sprintf("union(%s)", u.options.sorted().name())
}

func (u *union) Options() []types.Type {
	return u.options
}

func (u *union) Accepts(other types.Type) bool {
	if u == other {
		return true
	}
	if other, ok := other.(UnionType); ok {
		return u.acceptsUnion(other)
	}
	return u.acceptsType(other)
}

func (u *union) acceptsUnion(other UnionType) bool {
	for _, o := range other.Options() {
		if !u.acceptsType(o) {
			return false
		}
	}
	return true
}

func (u *union) acceptsType(other types.Type) bool {
	for _, t := range u.options {
		if t.Accepts(other) {
			return true
		}
	}
	return false
}

func basicType(o typeList) types.Type {
	first := o[0]
	for _, next := range o[1:] {
		if !first.Accepts(next) {
			return basic.New(o.name())
		}
	}
	return first
}
