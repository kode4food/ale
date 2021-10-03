package encoder_test

import (
	"testing"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
)

func TestCall(t *testing.T) {
	as := assert.New(t)
	f1 := func(_ encoder.Encoder, _ ...data.Value) {}
	c1 := encoder.Call(f1)
	as.String("encoder", c1.Type().Name())
	as.Contains(`:type encoder`, c1)
	as.Contains(`:instance `, c1)
}
