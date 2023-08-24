package generate

import "github.com/kode4food/ale/compiler/encoder"

// Builder is a callback that many of the functions in the package will invoke
// encode in-place instructions
type Builder func(encoder.Encoder)
