package special

import (
	"errors"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/data"
)

// Error messages
const (
	ErrPatternsNotSupported = "patterns not yet supported"
)

// Pattern instantiates a matchable pattern
func Pattern(encoder.Encoder, ...data.Value) {
	panic(errors.New(ErrPatternsNotSupported))
}
