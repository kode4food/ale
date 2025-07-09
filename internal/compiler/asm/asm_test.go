package asm_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/internal/compiler/ir/analysis"
)

func TestAsmStackSizeError(t *testing.T) {
	as := assert.New(t)
	as.ErrorWith(`(asm pop)`, fmt.Errorf(analysis.ErrBadStackTermination, -2))
}
