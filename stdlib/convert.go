package stdlib

import (
	"bytes"
	"fmt"

	"gitlab.com/kode4food/ale/data"
)

// SequenceToList takes any sequence and converts it to a List
func SequenceToList(s data.Sequence) *data.List {
	switch typed := s.(type) {
	case *data.List:
		return typed
	case data.Counted:
		res := make(data.Vector, typed.Count())
		idx := 0
		for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
			res[idx] = f
			idx++
		}
		return data.NewList(res...)
	default:
		return uncountedToList(typed)
	}
}

func uncountedToList(s data.Sequence) *data.List {
	return data.NewList(uncountedToVector(s)...)
}

// SequenceToValues takes any sequence and converts it to a value array
func SequenceToValues(s data.Sequence) data.Values {
	return data.Values(SequenceToVector(s))
}

// SequenceToVector takes any sequence and converts it to a vector
func SequenceToVector(s data.Sequence) data.Vector {
	switch typed := s.(type) {
	case data.Vector:
		return typed
	case data.Counted:
		res := make(data.Vector, typed.Count())
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

func uncountedToVector(s data.Sequence) data.Vector {
	res := data.Vector{}
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		res = append(res, f)
	}
	return res
}

// SequenceToAssociative takes any sequence and converts it to an Associative
func SequenceToAssociative(s data.Sequence) data.Associative {
	switch typed := s.(type) {
	case data.Associative:
		return typed
	case data.Counted:
		l := typed.Count()
		if l%2 != 0 {
			panic(fmt.Errorf(data.ExpectedPair))
		}
		ml := l / 2
		r := make([]data.Vector, ml)
		i := s
		var k, v data.Value
		for idx := 0; idx < ml; idx++ {
			k, i, _ = i.Split()
			v, i, _ = i.Split()
			r[idx] = data.Vector{k, v}
		}
		return data.Associative(r)
	default:
		return uncountedToAssociative(s)
	}
}

func uncountedToAssociative(s data.Sequence) data.Associative {
	res := make([]data.Vector, 0)
	var v data.Value
	for k, r, ok := s.Split(); ok; k, r, ok = r.Split() {
		if v, r, ok = r.Split(); ok {
			res = append(res, data.Vector{k, v})
		} else {
			panic(fmt.Errorf(data.ExpectedPair))
		}
	}
	return data.Associative(res)
}

// SequenceToStr takes any sequence and attempts to convert it to a String
func SequenceToStr(s data.Sequence) data.String {
	if st, ok := s.(data.String); ok {
		return st
	}
	var buf bytes.Buffer
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		if f == data.Nil {
			continue
		}
		buf.WriteString(f.String())
	}
	return data.String(buf.String())
}
