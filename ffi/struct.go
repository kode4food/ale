package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/sequence"
)

type (
	structWrapper struct {
		typ    reflect.Type
		fields []*fieldWrapper
	}

	fieldWrapper struct {
		Wrapper
		kwd data.Keyword
		idx int
	}
)

// AleTag identifies the tag used to specify the kwd used when wrapping a
// struct as an Object
const AleTag = "ale"

func makeWrappedStruct(t reflect.Type) (Wrapper, error) {
	fLen := t.NumField()
	var fields []*fieldWrapper
	for i := range fLen {
		f := t.Field(i)
		if f.PkgPath != "" { // Not exported
			continue
		}
		k := getFieldKeyword(f)
		w, err := WrapType(f.Type)
		if err != nil {
			return nil, err
		}
		fields = append(fields, &fieldWrapper{
			Wrapper: w,
			kwd:     k,
			idx:     i,
		})
	}
	return &structWrapper{
		typ:    t,
		fields: fields,
	}, nil
}

func getFieldKeyword(f reflect.StructField) data.Keyword {
	tag, ok := f.Tag.Lookup(AleTag)
	if !ok {
		tag = f.Name
	}
	return data.Keyword(tag)
}

func (w *structWrapper) Wrap(c *Context, v reflect.Value) (ale.Value, error) {
	if !v.IsValid() {
		return data.Null, nil
	}
	out := make(data.Pairs, 0, len(w.fields))
	for _, w := range w.fields {
		v, err := w.Wrap(c, v.Field(w.idx))
		if err != nil {
			return data.Null, err
		}
		out = append(out, data.NewCons(w.kwd, v))
	}
	return data.NewObject(out...), nil
}

func (w *structWrapper) Unwrap(v ale.Value) (reflect.Value, error) {
	s, ok := v.(data.Sequence)
	if !ok {
		return _zero, errors.New(ErrValueMustBeSequence)
	}
	in, err := sequence.ToObject(s)
	if err != nil {
		return _zero, err
	}
	out := reflect.New(w.typ).Elem()
	for _, w := range w.fields {
		if v, ok := in.Get(w.kwd); ok {
			v, err := w.Unwrap(v)
			if err != nil {
				return _zero, err
			}
			out.Field(w.idx).Set(v)
		}
	}
	return out, nil
}
