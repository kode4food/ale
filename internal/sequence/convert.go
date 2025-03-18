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
	default:
		v := ToVector(s)
		return data.NewList(v...)
	}
}

// ToVector takes any sequence and converts it to a Vector
func ToVector(s data.Sequence) data.Vector {
	switch s := s.(type) {
	case data.Vector:
		return s
	case data.CountedSequence:
		return countedToVector(s)
	default:
		return uncountedToVector(s)
	}
}

func countedToVector(s data.CountedSequence) data.Vector {
	res := make(data.Vector, s.Count())
	idx := 0
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		res[idx] = f
		idx++
	}
	return res
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
