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

func makeWrappedStruct(t reflect.Type) (Wrapper, error) {
	fLen := t.NumField()
	fields := make(map[string]*fieldWrapper, fLen)
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
		fields[f.Name] = &fieldWrapper{
			Wrapper: w,
			Keyword: k,
		}
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

func (s *structWrapper) Wrap(c *Context, v reflect.Value) (data.Value, error) {
	if !v.IsValid() {
		return data.Nil, nil
	}
	out := make(data.Pairs, 0, len(s.fields))
	for k, w := range s.fields {
		v, err := w.Wrap(c, v.FieldByName(k))
		if err != nil {
			return nil, err
		}
		out = append(out, data.NewCons(w.Keyword, v))
	}
	return data.NewObject(out...), nil
}

func (s *structWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	in, err := sequence.ToObject(v.(data.Sequence))
	if err != nil {
		return _emptyValue, err
	}
	out := reflect.New(s.typ).Elem()
	for k, w := range s.fields {
		if v, ok := in.Get(w.Keyword); ok {
			v, err := w.Unwrap(v)
			if err != nil {
				return _emptyValue, err
			}
			out.FieldByName(k).Set(v)
		}
	}
	return out, nil
}
