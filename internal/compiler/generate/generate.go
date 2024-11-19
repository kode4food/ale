package generate

import "github.com/kode4food/ale/internal/compiler/encoder"

// Builder is a callback that many of the functions in the package will invoke
// to encode in-place instructions
type Builder func(encoder.Encoder) error
