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
	as.String("special", c1.Type().Name())
	as.Contains(`:type special`, c1)
	as.Contains(`:instance `, c1)
}
