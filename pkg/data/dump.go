package data

import (
	"cmp"
	"fmt"
	"strings"

	"github.com/kode4food/ale/internal/basics"
	"github.com/kode4food/ale/internal/lang"
)

type dumped struct{ Value }

// Standard Keys
const (
	CountKey    = Keyword("count")
	InstanceKey = Keyword("instance")
	NameKey     = Keyword("name")
	TypeKey     = Keyword("type")
)

var (
	dumpMap = map[Value]func(Value) (Value, bool){
		CountKey:    dumpCount,
		InstanceKey: dumpInstance,
		NameKey:     dumpName,
		TypeKey:     dumpType,
	}

	dumpKeys = basics.SortedKeysFunc(dumpMap, func(l, r Value) int {
		return cmp.Compare(l.(Keyword), r.(Keyword))
	})
)

// DumpMapped takes a Value and dumps out a bunch of info as a Mapped
func DumpMapped(v Value) Mapped {
	return dump(v)
}

// DumpString takes a Value and dumps out a bunch of info as a string
func DumpString(v Value) string {
	return dump(v).String()
}

func dump(v Value) dumped {
	return dumped{v}
}

func (d dumped) Get(key Value) (Value, bool) {
	if f, ok := dumpMap[key]; ok {
		return f(d.Value)
	}
	return Null, false
}

func (d dumped) String() string {
	var buf strings.Builder
	buf.WriteString(lang.ObjectStart)
	i := 0
	for _, k := range dumpKeys {
		v, ok := d.Get(k)
		if !ok {
			continue
		}
		if i > 0 {
			buf.WriteString(lang.Space)
		}
		buf.WriteString(ToQuotedString(k))
		buf.WriteString(lang.Space)
		buf.WriteString(ToQuotedString(v))
		i++
	}
	buf.WriteString(lang.ObjectEnd)
	return buf.String()
}

func dumpCount(v Value) (Value, bool) {
	if c, ok := v.(Counted); ok {
		return c.Count(), true
	}
	return Null, false
}

func dumpInstance(v Value) (Value, bool) {
	return String(fmt.Sprintf("%p", v)), true
}

func dumpName(v Value) (Value, bool) {
	if n, ok := v.(Named); ok {
		return n.Name(), true
	}
	return Null, false
}

func dumpType(v Value) (Value, bool) {
	if t, ok := v.(Typed); ok {
		return Local(t.Type().Name()), true
	}
	return Null, false
}
