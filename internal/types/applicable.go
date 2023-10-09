package types

import (
	"fmt"
	"strings"
)

type (
	// Applicable describes a Type that can be the target of a function
	// application. Such application may expose multiple Signatures to
	// describe the various ways that the type can be applied
	Applicable struct {
		basic
		signatures
	}

	// Signature describes an ApplicableType calling signature
	Signature struct {
		Params    []Type
		TakesRest bool
		Result    Type
	}

	signatures []Signature
)

// MakeApplicable declares an ApplicableType that will only allow an applicable
// value capable of the provided signature set
func MakeApplicable(first Signature, rest ...Signature) *Applicable {
	all := append(signatures{first}, rest...)
	return &Applicable{
		basic:      BasicLambda,
		signatures: all,
	}
}

func (a *Applicable) Signatures() []Signature {
	return a.signatures
}

func (a *Applicable) Name() string {
	return fmt.Sprintf("%s(%s)", a.basic.Name(), a.signatures.name())
}

func (a *Applicable) Accepts(c *Checker, other Type) bool {
	if a == other {
		return true
	}
	if other, ok := other.(*Applicable); ok {
		os := other.Signatures()
		for _, s := range a.signatures {
			if !s.acceptsFromSignatures(c, os) {
				return false
			}
		}
		return true
	}
	return false
}

func (a *Applicable) Equal(other Type) bool {
	if a == other {
		return true
	}
	if other, ok := other.(*Applicable); ok {
		return a.basic.Equal(other.basic) &&
			a.signatures.equal(other.signatures)
	}
	return false
}

func (s Signature) name() string {
	return fmt.Sprintf("%s->%s", s.argNames(), s.Result.Name())
}

func (s Signature) argNames() string {
	a := s.Params
	if !s.TakesRest {
		return typeList(a).name()
	}
	l := len(a)
	args := typeList(a[:l-1]).name()
	rest := a[l-1].Name()
	return fmt.Sprintf("%s.%s", args, rest)
}

func (s Signature) acceptsFromSignatures(
	c *Checker, other []Signature,
) bool {
	for _, o := range other {
		if s.accepts(c, o) {
			return true
		}
	}
	return false
}

func (s Signature) accepts(c *Checker, other Signature) bool {
	if !c.AcceptsChild(s.Result, other.Result) {
		return false
	}
	sa := s.Params
	oa := other.Params
	if len(sa) != len(oa) || s.TakesRest != other.TakesRest {
		return false
	}
	for i, a := range sa {
		if !c.AcceptsChild(a, oa[i]) {
			return false
		}
	}
	return true
}

func (s Signature) equal(other Signature) bool {
	if s.TakesRest != other.TakesRest || !s.Result.Equal(other.Result) {
		return false
	}
	if len(s.Params) != len(other.Params) {
		return false
	}
	for i, l := range s.Params {
		if !l.Equal(other.Params[i]) {
			return false
		}
	}
	return true
}

func (s signatures) equal(other signatures) bool {
	if len(s) != len(other) {
		return false
	}
	for i, l := range s {
		r := other[i]
		if !l.equal(r) {
			return false
		}
	}
	return true
}

func (s signatures) name() string {
	return strings.Join(s.names(), ",")
}

func (s signatures) names() []string {
	res := make([]string, len(s))
	for i, sig := range s {
		res[i] = sig.name()
	}
	return res
}