package vm

import (
	"fmt"

	"github.com/kode4food/ale/data"
)

// Ref encapsulates a reference to a Value
type Ref struct {
	data.Value
}

func (r *Ref) String() string {
	return fmt.Sprintf("(ref %s)", r.Value.String())
}
