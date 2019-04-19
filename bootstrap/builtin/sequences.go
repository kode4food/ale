package builtin

import "gitlab.com/kode4food/ale/data"

func fetchSequence(args data.Vector) data.Sequence {
	return args[0].(data.Sequence)
}

// Seq attempts to convert the provided value to a sequence, or returns nil
func Seq(args ...data.Value) data.Value {
	if s, ok := args[0].(data.Sequence); ok {
		return s
	}
	return data.Nil
}

// First returns the first value in the sequence
func First(args ...data.Value) data.Value {
	return fetchSequence(args).First()
}

// Rest returns the sequence elements after the first value
func Rest(args ...data.Value) data.Value {
	return fetchSequence(args).Rest()
}

// Last returns the final element of the sequence
func Last(args ...data.Value) data.Value {
	var l data.Value
	for f, r, ok := fetchSequence(args).Split(); ok; f, r, ok = r.Split() {
		l = f
	}
	return l
}

// Cons prepends a value to the provided sequence
func Cons(args ...data.Value) data.Value {
	h := args[0]
	r := args[1]
	return r.(data.Sequence).Prepend(h)
}

// Conj conjoins a value to the provided sequence in some way
func Conj(args ...data.Value) data.Value {
	s := args[0].(data.Conjoiner)
	for _, f := range args[1:] {
		s = s.Conjoin(f).(data.Conjoiner)
	}
	return s
}

// Len returns the size of the provided sequence
func Len(args ...data.Value) data.Value {
	s := fetchSequence(args)
	l := data.Count(s)
	return data.Integer(l)
}

// Nth returns the nth element of the provided sequence
func Nth(args ...data.Value) data.Value {
	s := args[0].(data.Indexed)
	res, _ := s.ElementAt(int(args[1].(data.Integer)))
	return res
}

// Get returns a value by key from the provided mapped sequence
func Get(args ...data.Value) data.Value {
	s := args[0].(data.Mapped)
	res, _ := s.Get(args[1])
	return res
}

// IsSeq returns whether or not the provided value is a non-empty sequence
func IsSeq(args ...data.Value) data.Value {
	s, ok := args[0].(data.Sequence)
	return data.Bool(ok && s.IsSequence())
}

// IsLen returns whether or not the provided value is a countable sequence
func IsLen(args ...data.Value) data.Value {
	_, ok := args[0].(data.CountedSequence)
	return data.Bool(ok)
}

// IsIndexed returns whether or not the provided value is an indexed sequence
func IsIndexed(args ...data.Value) data.Value {
	_, ok := args[0].(data.IndexedSequence)
	return data.Bool(ok)
}
