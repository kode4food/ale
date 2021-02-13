package builtin

import (
	"errors"

	"github.com/kode4food/ale/data"
)

// Error messages
const (
	ErrIndexOutOfBounds = "index out of bounds"
	ErrPutRequiresPair  = "put requires a key/value combination or a pair"
)

// First returns the first value in the sequence
var First = data.Applicative(func(args ...data.Value) data.Value {
	return args[0].(data.Sequence).First()
}, 1)

// Rest returns the sequence elements after the first value
var Rest = data.Applicative(func(args ...data.Value) data.Value {
	return args[0].(data.Sequence).Rest()
}, 1)

// Append adds a value to the end of the provided AppenderSequence
var Append = data.Applicative(func(args ...data.Value) data.Value {
	a := args[0].(data.AppenderSequence)
	s := args[1]
	return a.Append(s)
}, 2)

// Reverse returns a reversed copy of a Sequence
var Reverse = data.Applicative(func(args ...data.Value) data.Value {
	r := args[0].(data.ReverserSequence)
	return r.Reverse()
}, 1)

// Length returns the element count of the provided sequence
var Length = data.Applicative(func(args ...data.Value) data.Value {
	s := args[0].(data.CountedSequence)
	l := s.Count()
	return data.Integer(l)
}, 1)

// Nth returns the nth element of the provided sequence or a default
var Nth = data.Applicative(func(args ...data.Value) data.Value {
	s := args[0].(data.IndexedSequence)
	i := int(args[1].(data.Integer))
	if res, ok := s.ElementAt(i); ok {
		return res
	}
	if len(args) > 2 {
		return args[2]
	}
	panic(errors.New(ErrIndexOutOfBounds))
}, 2, 3)

// IsSeq returns whether the provided value is a sequence
var IsSeq = data.Applicative(func(args ...data.Value) data.Value {
	_, ok := args[0].(data.Sequence)
	return data.Bool(ok)
}, 1)

// IsEmpty returns whether the provided sequence is empty
var IsEmpty = data.Applicative(func(args ...data.Value) data.Value {
	s := args[0].(data.Sequence)
	return data.Bool(s.IsEmpty())
}, 1)

// IsCounted returns whether the provided value is a counted sequence
var IsCounted = data.Applicative(func(args ...data.Value) data.Value {
	_, ok := args[0].(data.CountedSequence)
	return data.Bool(ok)
}, 1)

// IsIndexed returns whether the provided value is an indexed sequence
var IsIndexed = data.Applicative(func(args ...data.Value) data.Value {
	_, ok := args[0].(data.IndexedSequence)
	return data.Bool(ok)
}, 1)

// IsReverser returns whether the value is a reversible sequence
var IsReverser = data.Applicative(func(args ...data.Value) data.Value {
	_, ok := args[0].(data.ReverserSequence)
	return data.Bool(ok)
}, 1)
