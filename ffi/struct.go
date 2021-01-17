package ffi

import (
	"reflect"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/sequence"
)

type structWrapper struct {
	typ    reflect.Type
	fields map[string]Wrapper
}

func makeWrappedStruct(t reflect.Type) Wrapper {
	fLen := t.NumField()
	fields := make(map[string]Wrapper, fLen)
	for i := 0; i < fLen; i++ {
		f := t.Field(i)
		fields[f.Name] = wrapType(f.Type)
	}
	return &structWrapper{
		typ:    t,
		fields: fields,
	}
}

func (s *structWrapper) Wrap(v reflect.Value) data.Value {
	out := make(data.Object, len(s.fields))
	for k, w := range s.fields {
		out[data.Name(k)] = w.Wrap(v.FieldByName(k))
	}
	return out
}

func (s *structWrapper) Unwrap(v data.Value) reflect.Value {
	in := sequence.ToObject(v.(data.Sequence))
	out := reflect.New(s.typ).Elem()
	for k, w := range s.fields {
		v := w.Unwrap(in[data.Name(k)])
		out.FieldByName(k).Set(v)
	}
	return out
}
