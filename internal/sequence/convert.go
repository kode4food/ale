package sequence

import (
	"strings"

	"github.com/kode4food/ale/pkg/data"
)

// ToList takes any sequence and converts it to a List
func ToList(s data.Sequence) *data.List {
	switch s := s.(type) {
	case *data.List:
		return s
	case data.CountedSequence:
		res := make(data.Vector, s.Count())
		idx := 0
		for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
			res[idx] = f
			idx++
		}
		return data.NewList(res...)
	default:
		return uncountedToList(s)
	}
}

func uncountedToList(s data.Sequence) *data.List {
	return data.NewList(uncountedToVector(s)...)
}

// ToVector takes any sequence and converts it to a Vector
func ToVector(s data.Sequence) data.Vector {
	switch s := s.(type) {
	case data.Vector:
		return s
	case data.CountedSequence:
		res := make(data.Vector, s.Count())
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

// ToObject takes any sequence and converts it to an Object
func ToObject(s data.Sequence) (*data.Object, error) {
	switch s := s.(type) {
	case *data.Object:
		return s, nil
	default:
		v := ToVector(s)
		return data.ValuesToObject(v...)
	}
}

// ToString takes any sequence and attempts to convert it to a String
func ToString(s data.Sequence) data.String {
	if st, ok := s.(data.String); ok {
		return st
	}
	var buf strings.Builder
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		if f == data.Null {
			continue
		}
		buf.WriteString(data.ToString(f))
	}
	return data.String(buf.String())
}
