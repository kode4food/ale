package data

import "fmt"

// DumpString takes a Value and attempts to spit out a bunch of info
func DumpString(v Value) string {
	p := String(fmt.Sprintf("%p", v))
	m := ValueMap{InstanceKey: p}
	if n, ok := v.(Named); ok {
		m[NameKey] = n.Name()
	}
	if t, ok := v.(Typed); ok {
		m[TypeKey] = t.Type()
	}
	if c, ok := v.(Counted); ok {
		m[CountKey] = Integer(c.Count())
	}
	return m.String()
}
