package data

import (
	"bytes"
	"cmp"
	"fmt"

	"github.com/kode4food/comb/maps"
)

type dumpStringMap map[Value]Value

// DumpString takes a Value and attempts to dump out a bunch of info
func DumpString(v Value) string {
	p := String(fmt.Sprintf("%p", v))
	m := dumpStringMap{InstanceKey: p}
	if n, ok := v.(Named); ok {
		m[NameKey] = n.Name()
	}
	if t, ok := v.(Typed); ok {
		m[TypeKey] = Local(t.Type().Name())
	}
	if c, ok := v.(Counted); ok {
		m[CountKey] = Integer(c.Count())
	}
	return m.String()
}

var valueKeySorter = maps.SortedKeysFunc[Value, Value](func(l, r Value) int {
	return cmp.Compare(
		fmt.Sprintf("%p", l), fmt.Sprintf("%p", r),
	)
}).Must()

func (d dumpStringMap) sortedKeys() Values {
	return valueKeySorter(d)
}

func (d dumpStringMap) String() string {
	var buf bytes.Buffer
	buf.WriteString("{")
	for i, k := range d.sortedKeys() {
		if i > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(MaybeQuoteString(k))
		buf.WriteString(" ")
		buf.WriteString(MaybeQuoteString(d[k]))
	}
	buf.WriteString("}")
	return buf.String()
}
