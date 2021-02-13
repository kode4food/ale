package data

import (
	"bytes"
	"fmt"
	"sort"
)

// ValueMap is used to bootstrap a new Object
type ValueMap map[Value]Value

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
