package macro_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestQuoteObject(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(
		"(let [x :hello] `{,x 99})",
		O(data.NewCons(K("hello"), I(99))),
	)
}
