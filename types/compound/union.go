package compound

import (
	"bytes"

	"github.com/kode4food/ale/types"
)

type (
	// UnionType describes a Type that can accept any of a set of Types
	UnionType interface {
		types.Type
		union() // marker
		Options() []types.Type
	}

	union struct {
		options []types.Type
	}
)

// Union declares a UnionType based on at least one provided Type
func Union(first types.Type, rest ...types.Type) UnionType {
	all := append([]types.Type{first}, rest...)
	return &union{
		options: flattenUnions(all),
	}
}

func (*union) union() {}

func (u *union) Options() []types.Type {
	return u.options
}

func (u *union) Name() string {
	var buf bytes.Buffer
	buf.WriteString(u.options[0].Name())
	for _, n := range u.options[1:] {
		buf.WriteByte('|')
		buf.WriteString(n.Name())
	}
	return buf.String()
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

func flattenUnions(options []types.Type) []types.Type {
	var res []types.Type
	for _, o := range options {
		if o, ok := o.(UnionType); ok {
			res = append(res, o.Options()...)
			continue
		}
		res = append(res, o)
	}
	return res
}
