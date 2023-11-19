package builtin_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/core/builtin"
	"github.com/kode4food/ale/internal/assert"
)

func TestReachability(t *testing.T) {
	as := assert.New(t)

	as.PanicWith(`
		(lambda
			[(x y) "hello"]
			[(z) "there"]
			[(a b) "error"])`,
		fmt.Sprintf(builtin.ErrUnreachableCase, "(a b)", "(x y)"),
	)

	as.PanicWith(`
		(lambda
			[(x y . z) "hello"]
			[(x y) "there"]
			[(a b) "error"])`,
		fmt.Sprintf(builtin.ErrUnreachableCase, "(x y)", "(x y . z)"),
	)
}
