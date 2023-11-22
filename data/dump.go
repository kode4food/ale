package data

import (
	"bytes"
	"cmp"
	"fmt"

	"github.com/kode4food/comb/basics"
)

type dumped struct{ Value }

// Standard Keys
const (
	CountKey    = Keyword("count")
	HashKey     = Keyword("hash")
	InstanceKey = Keyword("instance")
	NameKey     = Keyword("name")
	TypeKey     = Keyword("type")
)

var (
	dumpMap = map[Value]func(Value) (Value, bool){
		CountKey:    dumpCount,
		HashKey:     dumpHash,
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
	var buf bytes.Buffer
	buf.WriteString("{")
	i := 0
	for _, k := range dumpKeys {
		v, ok := d.Get(k)
		if !ok {
			continue
		}
		if i > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(MaybeQuoteString(k))
		buf.WriteString(" ")
		buf.WriteString(MaybeQuoteString(v))
		i++
	}
	buf.WriteString("}")
	return buf.String()
}

func dumpCount(v Value) (Value, bool) {
	if c, ok := v.(Counted); ok {
		return Integer(c.Count()), true
	}
	return Null, false
}

func dumpHash(v Value) (Value, bool) {
	if h, ok := v.(Hashed); ok {
		return Integer(h.HashCode()), true
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
