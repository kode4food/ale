package ffi

import (
	"errors"
	"reflect"

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
		data.Keyword
		idx int
	}
)

// AleTag identifies the tag used to specify the Keyword used when wrapping a
// struct as an Object
const AleTag = "ale"

func makeWrappedStruct(t reflect.Type) (Wrapper, error) {
	fLen := t.NumField()
	var fields []*fieldWrapper
	for i := 0; i < fLen; i++ {
		f := t.Field(i)
		if f.PkgPath != "" { // Not exported
			continue
		}
		k := getFieldKeyword(f)
		w, err := wrapType(f.Type)
		if err != nil {
			return nil, err
		}
		fields = append(fields, &fieldWrapper{
			Wrapper: w,
			Keyword: k,
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

func (w *structWrapper) Wrap(c *Context, v reflect.Value) (data.Value, error) {
	if !v.IsValid() {
		return data.Nil, nil
	}
	out := make(data.Pairs, 0, len(w.fields))
	for _, w := range w.fields {
		v, err := w.Wrap(c, v.Field(w.idx))
		if err != nil {
			return data.Nil, err
		}
		out = append(out, data.NewCons(w.Keyword, v))
	}
	return data.NewObject(out...), nil
}

func (w *structWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	s, ok := v.(data.Sequence)
	if !ok {
		return _emptyValue, errors.New(ErrValueMustBeSequence)
	}
	in, err := sequence.ToObject(s)
	if err != nil {
		return _emptyValue, err
	}
	out := reflect.New(w.typ).Elem()
	for _, w := range w.fields {
		if v, ok := in.Get(w.Keyword); ok {
			v, err := w.Unwrap(v)
			if err != nil {
				return _emptyValue, err
			}
			out.Field(w.idx).Set(v)
		}
	}
	return out, nil
}
