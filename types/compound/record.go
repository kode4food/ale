package compound

import (
	"github.com/kode4food/ale/types"
	"github.com/kode4food/ale/types/basic"
)

type (
	// RecordType describes an Object that allows a fixed set of fields,
	// each of which has a keyword
	RecordType interface {
		types.BasicType
		record() // marker
		Fields() []Field
	}

	// Field describes one of the fields of a RecordType
	Field struct {
		Name  string
		Value types.Type
	}

	record struct {
		types.BasicType
		fields []Field
	}
)

// Record declares a new RecordType that only allows a fixed set of Field
// entries, each being identified by a Keyword and having a specified Type
func Record(fields ...Field) RecordType {
	return &record{
		BasicType: basic.Object,
		fields:    fields,
	}
}

func (*record) record() {}

func (r *record) Fields() []Field {
	return r.fields
}

func (*record) Name() string {
	return "record"
}

func (r *record) Accepts(other types.Type) bool {
	if r == other {
		return true
	}
	if other, ok := other.(RecordType); ok {
		rf := r.fields
		of := other.Fields()
		if len(rf) > len(of) {
			return false
		}
		om := fieldsToMap(of)
		for k, v := range fieldsToMap(rf) {
			if tv, ok := om[k]; !ok || !v.Accepts(tv) {
				return false
			}
		}
		return true
	}
	return false
}

func fieldsToMap(fields []Field) map[string]types.Type {
	res := map[string]types.Type{}
	for _, f := range fields {
		res[f.Name] = f.Value
	}
	return res
}
