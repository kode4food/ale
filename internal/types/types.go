package types

import (
	"cmp"
	"slices"
	"strings"

	"github.com/kode4food/comb/basics"
)

type (
	// Type describes the type compatibility for a Value
	Type interface {
		// Name identifies this Type
		Name() string

		// Accepts determines if this Type will accept the provided Type for
		// binding. This will generally mean that the provided Type satisfies
		// the contract of the receiver. A Checker is provided to track the
		// state of the Type checking
		Accepts(*Checker, Type) bool

		// Equal determines if the provided Type is an equivalent definition
		Equal(Type) bool
	}

	typeList []Type
)

func Equal(l, r Type) bool {
	return l.Equal(r)
}

func (t typeList) sorted() typeList {
	return basics.SortFunc(t, func(l, r Type) int {
		return cmp.Compare(l.Name(), r.Name())
	})
}

func (t typeList) deduplicated() typeList {
	var res typeList
	var last Type
	for _, t := range t.sorted() {
		if t == last {
			continue
		}
		res = append(res, t)
		last = t
	}
	return res
}

func (t typeList) name() string {
	return strings.Join(t.names(), ",")
}

func (t typeList) names() []string {
	res := make([]string, len(t))
	for i, t := range t {
		res[i] = t.Name()
	}
	return res
}

func (t typeList) flatten() typeList {
	var res typeList
	for _, o := range t {
		if o, ok := o.(*Union); ok {
			res = append(res, o.Options()...)
			continue
		}
		res = append(res, o)
	}
	return res.deduplicated()
}

func (t typeList) hasAny() bool {
	for _, t := range t {
		if _, ok := t.(*Any); ok {
			return true
		}
	}
	return false
}

func (t typeList) basicType() (basic, bool) {
	f, ok := t[0].(basic)
	if !ok {
		return nil, false
	}
	for _, n := range t[1:] {
		if n, ok := n.(basic); ok && f.Kind() == n.Kind() {
			continue
		}
		return nil, false
	}
	return f, true
}

func (t typeList) equal(other typeList) bool {
	tf := t.flatten()
	of := other.flatten()
	return slices.EqualFunc(tf, of, func(l, r Type) bool {
		return l.Equal(r)
	})
}
