package data

import (
	"bytes"
	"cmp"
	"fmt"
	"slices"
)

type dumpStringMap map[Value]Value

// DumpString takes a Value and attempts to spit out a bunch of info
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

func (d dumpStringMap) sortedKeys() Values {
	keys := make(Values, 0, len(d))
	for k := range d {
		keys = append(keys, k)
	}
	slices.SortFunc(keys, func(l, r Value) int {
		return cmp.Compare(
			fmt.Sprintf("%p", l), fmt.Sprintf("%p", r),
		)
	})
	return keys
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
