package stdlib

import (
	"bytes"
	"fmt"

	"gitlab.com/kode4food/ale/api"
)

// SequenceToList takes any sequence and converts it to a List
func SequenceToList(s api.Sequence) *api.List {
	switch typed := s.(type) {
	case *api.List:
		return typed
	case api.Counted:
		res := make(api.Vector, typed.Count())
		idx := 0
		for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
			res[idx] = f
			idx++
		}
		return api.NewList(res...)
	default:
		return uncountedToList(typed)
	}
}

func uncountedToList(s api.Sequence) *api.List {
	return api.NewList(uncountedToVector(s)...)
}

// SequenceToValues takes any sequence and converts it to a value array
func SequenceToValues(s api.Sequence) api.Values {
	return api.Values(SequenceToVector(s))
}

// SequenceToVector takes any sequence and converts it to a vector
func SequenceToVector(s api.Sequence) api.Vector {
	switch typed := s.(type) {
	case api.Vector:
		return typed
	case api.Counted:
		res := make(api.Vector, typed.Count())
		idx := 0
		for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
			res[idx] = f
			idx++
		}
		return res
	default:
		return uncountedToVector(s)
	}
}

func uncountedToVector(s api.Sequence) api.Vector {
	res := api.Vector{}
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		res = append(res, f)
	}
	return res
}

// SequenceToAssociative takes any sequence and converts it to an Associative
func SequenceToAssociative(s api.Sequence) api.Associative {
	switch typed := s.(type) {
	case api.Associative:
		return typed
	case api.Counted:
		l := typed.Count()
		if l%2 != 0 {
			panic(fmt.Errorf(api.ExpectedPair))
		}
		ml := l / 2
		r := make([]api.Vector, ml)
		i := s
		var k, v api.Value
		for idx := 0; idx < ml; idx++ {
			k, i, _ = i.Split()
			v, i, _ = i.Split()
			r[idx] = api.Vector{k, v}
		}
		return api.Associative(r)
	default:
		return uncountedToAssociative(s)
	}
}

func uncountedToAssociative(s api.Sequence) api.Associative {
	res := make([]api.Vector, 0)
	var v api.Value
	for k, r, ok := s.Split(); ok; k, r, ok = r.Split() {
		if v, r, ok = r.Split(); ok {
			res = append(res, api.Vector{k, v})
		} else {
			panic(fmt.Errorf(api.ExpectedPair))
		}
	}
	return api.Associative(res)
}

// SequenceToStr takes any sequence and attempts to convert it to a String
func SequenceToStr(s api.Sequence) api.String {
	if st, ok := s.(api.String); ok {
		return st
	}
	var buf bytes.Buffer
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		if f == api.Nil {
			continue
		}
		buf.WriteString(f.String())
	}
	return api.String(buf.String())
}
