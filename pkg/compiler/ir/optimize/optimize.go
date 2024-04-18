package optimize

import "github.com/kode4food/ale/pkg/compiler/encoder"

type optimizer func(*encoder.Encoded)

var optimizers = []optimizer{
	splitReturns,   // roll standalone returns into preceding branches
	literalReturns, // convert literal returns into single instructions
	makeTailCalls,  // replace calls in tail position with a tail-call
	inlineCalls,    // inline calls to procedures that qualify
}

// Encoded takes an Encoded representation and returns an optimized one
func Encoded(e *encoder.Encoded) *encoder.Encoded {
	res := e.Copy()
	for _, o := range optimizers {
		o(res)
	}
	return res
}
