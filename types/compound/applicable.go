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
		types.Type
		applicable() // marker
		Signatures() []Signature
	}

	// Signature describes an ApplicableType calling signature
	Signature struct {
		Arguments []types.Type
		Result    types.Type
	}

	applicable struct {
		types.Type
		signatures []Signature
	}
)

// Applicable declares an ApplicableType that will only allow an applicable
// value capable of the provided Signature set
func Applicable(s ...Signature) ApplicableType {
	return &applicable{
		Type:       basic.Lambda,
		signatures: s,
	}
}

func (a *applicable) applicable() {}

func (a *applicable) Signatures() []Signature {
	return a.signatures
}

func (a *applicable) Accepts(t types.Type) bool {
	panic("implement me")
}
