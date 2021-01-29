package ffi

import (
	"reflect"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/sequence"
)

type (
	structWrapper struct {
		typ    reflect.Type
		fields map[string]*fieldWrapper
	}

	fieldWrapper struct {
		Wrapper
		data.Keyword
	}
)

// AleTag identifies the tag used to specify the Keyword used when
// wrapping a struct as an Object
const AleTag = "ale"

func makeWrappedStruct(t reflect.Type) Wrapper {
	fLen := t.NumField()
	fields := make(map[string]*fieldWrapper, fLen)
	for i := 0; i < fLen; i++ {
		f := t.Field(i)
		if f.PkgPath != "" { // Not exported
			continue
		}
		k := getFieldKeyword(f)
		fields[f.Name] = &fieldWrapper{
			Wrapper: wrapType(f.Type),
			Keyword: k,
		}
	}
	return &structWrapper{
		typ:    t,
		fields: fields,
	}
}

func getFieldKeyword(f reflect.StructField) data.Keyword {
	tag, ok := f.Tag.Lookup(AleTag)
	if !ok {
		tag = f.Name
	}
	return data.Keyword(tag)
}

func (s *structWrapper) Wrap(c *WrapContext, v reflect.Value) data.Value {
	if r, ok := c.Get(v); ok {
		return r
	}
	out := make(data.Object, len(s.fields))
	c.Put(v, out)
	for k, w := range s.fields {
		out[w.Keyword] = w.Wrap(c, v.FieldByName(k))
	}
	return out
}

func (s *structWrapper) Unwrap(c *UnwrapContext, v data.Value) reflect.Value {
	if r, ok := c.Get(v); ok {
		return r
	}
	in := sequence.ToObject(v.(data.Sequence))
	out := reflect.New(s.typ).Elem()
	c.Put(v, out)
	for k, w := range s.fields {
		if v, ok := in[w.Keyword]; ok {
			out.FieldByName(k).Set(w.Unwrap(c, v))
		}
	}
	return out
}
