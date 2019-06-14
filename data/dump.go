package data

import "fmt"

// DumpString takes a Value and attempts to spit out a bunch of info
func DumpString(v Value) string {
	p := String(fmt.Sprintf("%p", v))
	m := Object{InstanceKey: p}
	if t, ok := v.(Typed); ok {
		m = m.Copy()
		m[TypeKey] = t.Type()
	}
	return m.String()
}
