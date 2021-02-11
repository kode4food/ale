package data

import (
	"bytes"
	"fmt"
	"sort"
)

// ValueMap is used to bootstrap a new Object
type ValueMap map[Value]Value

func (v ValueMap) firstFrom(keys []Value) Value {
	if len(keys) > 0 {
		k0 := keys[0]
		return NewCons(k0, v[k0])
	}
	return Nil
}

func (v ValueMap) restFrom(keys []Value) ValueMap {
	if len(keys) > 1 {
		rest := make(ValueMap, len(keys)-1)
		for _, k := range keys[1:] {
			rest[k] = v[k]
		}
		return rest
	}
	return ValueMap{}
}

func (v ValueMap) sortedKeys() []Value {
	keys := make([]Value, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(l, r int) bool {
		return fmt.Sprintf("%p", keys[l]) < fmt.Sprintf("%p", keys[r])
	})
	return keys
}

// String converts this ValueMap into a string
func (v ValueMap) String() string {
	var buf bytes.Buffer
	buf.WriteString("{")
	for i, k := range v.sortedKeys() {
		if i > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(MaybeQuoteString(k))
		buf.WriteString(" ")
		buf.WriteString(MaybeQuoteString(v[k]))
	}
	buf.WriteString("}")
	return buf.String()
}
