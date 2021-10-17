package compound

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/kode4food/ale/types"
	"github.com/kode4food/ale/types/basic"
	"github.com/kode4food/ale/types/extended"
)

type (
	// RecordType describes an Object that allows a fixed set of fields,
	// each of which has a keyword
	RecordType interface {
		types.Extended
		record() // marker
		Fields() []Field
	}

	// Field describes one of the fields of a RecordType
	Field struct {
		Name  string
		Value types.Type
	}

	fields []Field

	record struct {
		types.Extended
		fields
	}
)

// Record declares a new RecordType that only allows a fixed set of Field
// entries, each being identified by a Keyword and having a specified Type
func Record(fields ...Field) RecordType {
	return &record{
		Extended: extended.New(basic.Object),
		fields:   fields,
	}
}

func (*record) record() {}

func (r *record) Fields() []Field {
	return r.fields
}

func (r *record) Name() string {
	return fmt.Sprintf("record(%s)", r.fields.name())
}

func (r *record) Accepts(c types.Checker, other types.Type) bool {
	if r == other {
		return true
	}
	if other, ok := other.(RecordType); ok {
		rf := r.fields
		of := other.Fields()
		if len(rf) > len(of) {
			return false
		}
		om := fields(of).toMap()
		for k, v := range rf.toMap() {
			if tv, ok := om[k]; !ok || c.Check(v).Accepts(tv) == nil {
				return false
			}
		}
		return true
	}
	return false
}

func (f fields) toMap() map[string]types.Type {
	res := map[string]types.Type{}
	for _, p := range f {
		res[p.Name] = p.Value
	}
	return res
}

func (f fields) name() string {
	in := f[:]
	sort.Slice(in, func(i, j int) bool {
		return in[i].Name < in[j].Name
	})
	var buf bytes.Buffer
	for i, elem := range in {
		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteRune('"')
		buf.WriteString(elem.Name) // todo: escape the value
		buf.WriteString("\"->")
		buf.WriteString(elem.Value.Name())
	}
	return buf.String()
}
