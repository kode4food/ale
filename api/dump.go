package api

import "fmt"

const (
	instanceKey = Keyword("instance")
	typeKey     = Keyword("type")
)

// DumpString takes a Value and attempts to spit out a bunch of info
func DumpString(v Value) string {
	p := String(fmt.Sprintf("%p", v))
	m := Object{instanceKey: p}
	if t, ok := v.(Typed); ok {
		m = m.Extend(Object{typeKey: t.Type()})
	}
	return m.String()
}
