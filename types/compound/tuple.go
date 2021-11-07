package compound

import (
	"fmt"

	"github.com/kode4food/ale/types"
	"github.com/kode4food/ale/types/basic"
	"github.com/kode4food/ale/types/extended"
)

// Tuple declares a new TupleType that will only allow a List or Vector with
// positional elements of the provided Types
func Tuple(elems ...types.Type) types.Type {
	return makeTuple(basic.List, elems)
}

func makeTuple(base types.Type, elems []types.Type) types.Type {
	e := extended.New(base)
	var res types.Type = basic.Null
	for i := len(elems) - 1; i >= 0; i = i - 1 {
		res = &namedPair{
			name: fmt.Sprintf("tuple(%s)", typeList(elems).name()),
			pair: &pair{
				Extended: e,
				car:      elems[i],
				cdr:      res,
			},
		}
	}
	return res
}
