package stdlib

import (
	"bytes"

	"gitlab.com/kode4food/ale/data"
)

// SequenceToList takes any sequence and converts it to a List
func SequenceToList(s data.Sequence) data.List {
	switch typed := s.(type) {
	case data.List:
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

func uncountedToList(s data.Sequence) data.List {
	return data.NewList(uncountedToValues(s)...)
}

// SequenceToValues takes any sequence and converts it to a value array
func SequenceToValues(s data.Sequence) data.Values {
	switch typed := s.(type) {
	case data.Vector:
		return data.Values(typed)
	case data.Counted:
		res := make(data.Values, typed.Count())
		idx := 0
		for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
			res[idx] = f
			idx++
		}
		return res
	default:
		return uncountedToValues(s)
	}
}

func uncountedToValues(s data.Sequence) data.Values {
	res := data.Values{}
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		res = append(res, f)
	}
	return res
}

// SequenceToVector takes any sequence and converts it to a vector
func SequenceToVector(s data.Sequence) data.Vector {
	v := SequenceToValues(s)
	return data.NewVector(v...)
}

// SequenceToAssociative takes any sequence and converts it to an Associative
func SequenceToAssociative(s data.Sequence) data.Associative {
	switch typed := s.(type) {
	case data.Associative:
		return typed
	default:
		elems := SequenceToValues(s)
		return data.NewAssociative(elems...)
	}
}

// SequenceToStr takes any sequence and attempts to convert it to a String
func SequenceToStr(s data.Sequence) data.String {
	if st, ok := s.(data.String); ok {
		return st
	}
	var buf bytes.Buffer
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		if f == data.Null {
			continue
		}
		buf.WriteString(f.String())
	}
	return data.String(buf.String())
}
