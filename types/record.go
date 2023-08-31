package types

import (
	"bytes"
	"cmp"
	"fmt"
	"slices"
	"strings"
)

type (
	// RecordType describes an Object that allows a fixed set of fields,
	// each of which has a keyword
	RecordType interface {
		Type
		record() // marker
		Fields() []Field
	}

	// Field describes one of the fields of a RecordType
	Field struct {
		Name  string
		Value Type
	}

	fields []Field

	record struct {
		BasicType
		fields
	}
)

// Record declares a new RecordType that only allows a fixed set of Field
// entries, each being identified by a Keyword and having a specified Type
func Record(fields ...Field) RecordType {
	return &record{
		BasicType: AnyObject,
		fields:    fields,
	}
}

func (*record) record() {}

func (r *record) Fields() []Field {
	return r.fields
}

func (r *record) Name() string {
	return fmt.Sprintf("record(%s)", r.fields.name())
}

func (r *record) Accepts(c *Checker, other Type) bool {
	if other, ok := other.(RecordType); ok {
		rf := r.fields
		of := other.Fields()
		if len(rf) > len(of) {
			return false
		}
		om := fields(of).toMap()
		for k, v := range rf.toMap() {
			if tv, ok := om[k]; !ok || !c.AcceptsChild(v, tv) {
				return false
			}
		}
		return true
	}
	return false
}

func (f fields) toMap() map[string]Type {
	res := map[string]Type{}
	for _, p := range f {
		res[p.Name] = p.Value
	}
	return res
}

func (f fields) sorted() fields {
	res := make(fields, len(f))
	copy(res, f)
	slices.SortFunc(res, func(l, r Field) int {
		return cmp.Compare(l.Name, r.Name)
	})
	return res
}

func (f fields) name() string {
	var buf bytes.Buffer
	for i, elem := range f.sorted() {
		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteRune('"')
		buf.WriteString(strings.ReplaceAll(elem.Name, `"`, `\"`))
		buf.WriteString("\"->")
		buf.WriteString(elem.Value.Name())
	}
	return buf.String()
}
