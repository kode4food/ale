package types

import (
	"cmp"
	"slices"
	"strings"
)

type (
	// Type describes the type compatibility for a Value
	Type interface {
		// Name identifies this Type
		Name() string

		IsA(BasicType) bool

		// Accepts determines if this Type will accept the provided Type for
		// binding. This will generally mean that the provided Type satisfies
		// the contract of the receiver. A Checker is provided that tracks
		// the state of the Type checking
		Accepts(*Checker, Type) bool
	}

	typeList []Type
)

func (t typeList) sorted() typeList {
	res := make(typeList, len(t))
	copy(res, t)
	slices.SortFunc(res, func(l, r Type) int {
		return cmp.Compare(l.Name(), r.Name())
	})
	return res
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
		if o, ok := o.(UnionType); ok {
			res = append(res, o.Options()...)
			continue
		}
		res = append(res, o)
	}
	return res.deduplicated()
}

func (t typeList) hasAny() bool {
	for _, t := range t {
		if _, ok := t.(*anyType); ok {
			return true
		}
	}
	return false
}

func (t typeList) basicType() (BasicType, bool) {
	f, ok := t[0].(BasicType)
	if !ok {
		return nil, false
	}
	for _, n := range t[1:] {
		if !n.IsA(f) {
			return nil, false
		}
	}
	return f, true
}
