package compiler_test

import (
	"testing"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/internal/compiler"
	"github.com/kode4food/ale/internal/compiler/encoder"
)

func TestCall(t *testing.T) {
	as := assert.New(t)
	f1 := func(encoder.Encoder, ...ale.Value) error { return nil }
	c1 := compiler.Call(f1)
	as.True(compiler.CallType.Accepts(c1.Type()))
	as.False(c1.Type().Accepts(compiler.CallType))
	as.Contains(`special(0x`, c1)
}
