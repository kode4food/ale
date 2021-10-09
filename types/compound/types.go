package compound

import (
	"sort"
	"strings"

	"github.com/kode4food/ale/types"
	"github.com/kode4food/ale/types/basic"
)

type typeList []types.Type

func (t typeList) sorted() typeList {
	res := t[:]
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name() < res[j].Name()
	})
	return res
}

func (t typeList) deduplicated() typeList {
	var res typeList
	var last types.Type
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
		if _, ok := t.(basic.AnyType); ok {
			return true
		}
	}
	return false
}
