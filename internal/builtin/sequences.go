package builtin

import "gitlab.com/kode4food/ale/api"

func fetchSequence(args api.Vector) api.Sequence {
	return args[0].(api.Sequence)
}

// Seq attempts to convert the provided value to a sequence, or returns nil
func Seq(args ...api.Value) api.Value {
	if s, ok := args[0].(api.Sequence); ok {
		return s
	}
	return api.Nil
}

// First returns the first value in the sequence
func First(args ...api.Value) api.Value {
	return fetchSequence(args).First()
}

// Rest returns the sequence elements after the first value
func Rest(args ...api.Value) api.Value {
	return fetchSequence(args).Rest()
}

// Last returns the final element of the sequence
func Last(args ...api.Value) api.Value {
	var l api.Value
	for f, r, ok := fetchSequence(args).Split(); ok; f, r, ok = r.Split() {
		l = f
	}
	return l
}

// Cons prepends a value to the provided sequence
func Cons(args ...api.Value) api.Value {
	h := args[0]
	r := args[1]
	return r.(api.Sequence).Prepend(h)
}

// Conj conjoins a value to the provided sequence in some way
func Conj(args ...api.Value) api.Value {
	s := args[0].(api.Conjoiner)
	for _, f := range args[1:] {
		s = s.Conjoin(f).(api.Conjoiner)
	}
	return s
}

// Len returns the size of the provided sequence
func Len(args ...api.Value) api.Value {
	s := fetchSequence(args)
	l := api.Count(s)
	return api.Integer(l)
}

// Nth returns the nth element of the provided sequence
func Nth(args ...api.Value) api.Value {
	s := args[0].(api.Indexed)
	res, _ := s.ElementAt(int(args[1].(api.Integer)))
	return res
}

// Get returns a value by key from the provided mapped sequence
func Get(args ...api.Value) api.Value {
	s := args[0].(api.Mapped)
	res, _ := s.Get(args[1])
	return res
}

// IsSeq returns whether or not the provided value is a non-empty sequence
func IsSeq(args ...api.Value) api.Value {
	s, ok := args[0].(api.Sequence)
	return api.Bool(ok && s.IsSequence())
}

// IsLen returns whether or not the provided value is a countable sequence
func IsLen(args ...api.Value) api.Value {
	_, ok := args[0].(api.CountedSequence)
	return api.Bool(ok)
}

// IsIndexed returns whether or not the provided value is an indexed sequence
func IsIndexed(args ...api.Value) api.Value {
	_, ok := args[0].(api.IndexedSequence)
	return api.Bool(ok)
}
