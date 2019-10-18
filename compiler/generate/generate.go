package generate

import "github.com/kode4food/ale/data"

// TailCaller marks a Function as being capable of tail call
// optimized. Most specifically it is implemented by the
// Closure struct to inform the code generator
type TailCaller interface {
	data.Function
	TailCaller()
}
