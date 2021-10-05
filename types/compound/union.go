package compound

import (
	"bytes"
	"sort"

	"github.com/kode4food/ale/types"
)

type (
	// UnionType describes a Type that can accept any of a set of Types
	UnionType interface {
		types.BasicType
		union() // marker
		Options() Options
	}

	// Options stores the type.Type options of a UnionType
	Options []types.Type

	union struct {
		types.BasicType
		options Options
	}
)

// Union declares a UnionType based on at least one provided Type
func Union(first types.Type, rest ...types.Type) UnionType {
	all := append(Options{first}, rest...).flatten()
	return &union{
		BasicType: all.basicType(),
		options:   all,
	}
}

func (*union) union() {}

func (u *union) Name() string {
	return u.options.name()
}

func (u *union) Options() Options {
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

func (o Options) name() string {
	var buf bytes.Buffer
	s := o.sorted().names()
	buf.WriteString(s[0])
	for _, n := range s[1:] {
		buf.WriteByte('|')
		buf.WriteString(n)
	}
	return buf.String()
}

func (o Options) names() []string {
	res := make([]string, len(o))
	for i, t := range o {
		res[i] = t.Name()
	}
	return res
}

func (o Options) basicType() types.BasicType {
	if first, ok := o[0].(types.BasicType); ok {
		for _, next := range o[1:] {
			if !first.Accepts(next) {
				return types.Basic(o.name())
			}
		}
		return first
	}
	return types.Basic(o.name())
}

func (o Options) flatten() Options {
	var res Options
	for _, o := range o {
		if o, ok := o.(UnionType); ok {
			res = append(res, o.Options()...)
			continue
		}
		res = append(res, o)
	}
	return res.deduped()
}

func (o Options) sorted() Options {
	res := o[:]
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name() < res[j].Name()
	})
	return res
}

func (o Options) deduped() Options {
	var res Options
	var last types.Type
	for _, t := range o.sorted() {
		if t == last {
			continue
		}
		res = append(res, t)
		last = t
	}
	return res
}
