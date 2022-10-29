package types

import (
	"sort"
	"strings"
)

type (
	// Type describes the type compatibility for a Value
	Type interface {
		// Name identifies this Type
		Name() string

		IsA(BasicType) bool

		// Accepts determined if this Type will accept the provided Type for
		// binding. This will generally mean that the provided Type satisfies
		// the contract of the receiver. A Checker is provided that tracks
		// the state of the Type checking
		Accepts(*Checker, Type) bool
	}

	typeList []Type
)

func (t typeList) sorted() typeList {
	res := t[:]
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name() < res[j].Name()
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
	if f, ok := t[0].(BasicType); ok {
		for _, n := range t[1:] {
			if !n.IsA(f) {
				return nil, false
			}
		}
		return f, true
	}
	return nil, false
}
