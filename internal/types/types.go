package types

import (
	"cmp"
	"strings"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/internal/basics"
)

type typeList []ale.Type

func Equal(l, r ale.Type) bool {
	return l.Equal(r)
}

func (t typeList) sorted() typeList {
	return basics.SortedFunc(t, func(l, r ale.Type) int {
		return cmp.Compare(l.Name(), r.Name())
	})
}

func (t typeList) deduplicated() typeList {
	var res typeList
	var last ale.Type
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
	res := typeList{}
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
		if n, ok := n.(basic); ok && f.ID() == n.ID() {
			continue
		}
		return nil, false
	}
	return f, true
}

func (t typeList) equal(other typeList) bool {
	tf := t.flatten()
	of := other.flatten()
	return basics.EqualFunc(tf, of, func(l, r ale.Type) bool {
		return l.Equal(r)
	})
}
