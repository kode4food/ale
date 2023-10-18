package vm_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/runtime/vm"
)

func TestRef(t *testing.T) {
	as := assert.New(t)

	r1 := &vm.Ref{Value: S("hello")}
	r2 := &vm.Ref{Value: S("non-matching")}
	r3 := &vm.Ref{Value: S("hello")}
	r4 := &vm.Ref{Value: L(I(1), I(2))}
	r5 := &vm.Ref{Value: L(I(1), I(2))}
	r6 := &vm.Ref{}

	as.Equal(r1, r1)
	as.Equal(r1, r3)
	as.Equal(r1, S("hello"))
	as.NotEqual(r1, r2)
	as.Equal(r4, r5)
	as.Equal(r4, L(I(1), I(2)))
	as.NotEqual(r5, r6)

	as.String(`(ref "non-matching")`, r2)
	as.String(`(ref (1 2))`, r4)
	as.String(`(ref)`, r6)
}
