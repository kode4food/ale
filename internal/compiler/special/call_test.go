package special_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/special"
	"github.com/kode4food/ale/pkg/data"
)

func TestCall(t *testing.T) {
	as := assert.New(t)
	f1 := func(encoder.Encoder, ...data.Value) {}
	c1 := special.Call(f1)
	as.String("special", c1.Type().Name())
	as.Contains(`:type special`, c1)
	as.Contains(`:instance `, c1)
}
