package compound

import (
	"github.com/kode4food/ale/types"
	"github.com/kode4food/ale/types/basic"
)

type (
	// ApplicableType describes a Type that can be the target of a function
	// application. Such application may expose multiple Signatures to
	// describe the various ways that the type can be applied
	ApplicableType interface {
		types.BasicType
		applicable() // marker
		Signatures() []Signature
	}

	// Signature describes an ApplicableType calling signature
	Signature struct {
		Arguments []types.Type
		Result    types.Type
	}

	applicable struct {
		types.BasicType
		signatures []Signature
	}
)

// Applicable declares an ApplicableType that will only allow an applicable
// value capable of the provided Signature set
func Applicable(first Signature, rest ...Signature) ApplicableType {
	all := append([]Signature{first}, rest...)
	return &applicable{
		BasicType:  basic.Lambda,
		signatures: all,
	}
}

func (a *applicable) applicable() {}

func (a *applicable) Signatures() []Signature {
	return a.signatures
}

func (a *applicable) Accepts(other types.Type) bool {
	if a == other {
		return true
	}
	if other, ok := other.(ApplicableType); ok {
		os := other.Signatures()
		for _, s := range a.signatures {
			if !s.acceptsFromSignatures(os) {
				return false
			}
		}
		return true
	}
	return false
}

func (s Signature) acceptsFromSignatures(other []Signature) bool {
	for _, o := range other {
		if s.accepts(o) {
			return true
		}
	}
	return false
}

func (s Signature) accepts(other Signature) bool {
	if !s.Result.Accepts(other.Result) {
		return false
	}
	sa := s.Arguments
	oa := other.Arguments
	if len(sa) != len(oa) {
		return false
	}
	for i, a := range sa {
		if !a.Accepts(oa[i]) {
			return false
		}
	}
	return true
}
