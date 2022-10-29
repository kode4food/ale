package types

import (
	"fmt"
	"strings"
)

type (
	// ApplicableType describes a Type that can be the target of a function
	// application. Such application may expose multiple Signatures to
	// describe the various ways that the type can be applied
	ApplicableType interface {
		Type
		applicable() // marker
		Signatures() []Signature
		Accepts(*Checker, Type) bool
	}

	// Signature describes an ApplicableType calling signature
	Signature struct {
		Arguments []Type
		TakesRest bool
		Result    Type
	}

	applicable struct {
		BasicType
		signatures
	}

	signatures []Signature
)

// Applicable declares an ApplicableType that will only allow an applicable
// value capable of the provided Signature set
func Applicable(first Signature, rest ...Signature) ApplicableType {
	all := append(signatures{first}, rest...)
	return &applicable{
		BasicType:  Lambda,
		signatures: all,
	}
}

func (a *applicable) applicable() {}

func (a *applicable) Signatures() []Signature {
	return a.signatures
}

func (a *applicable) Name() string {
	return fmt.Sprintf("%s(%s)", a.BasicType.Name(), a.signatures.name())
}

func (a *applicable) Accepts(c *Checker, other Type) bool {
	if other, ok := other.(ApplicableType); ok {
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

func (s Signature) name() string {
	return fmt.Sprintf("%s->%s", s.argNames(), s.Result.Name())
}

func (s Signature) argNames() string {
	a := s.Arguments
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
	sa := s.Arguments
	oa := other.Arguments
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
