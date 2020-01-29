package arity_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/compiler/arity"
	"github.com/kode4food/ale/internal/assert"
)

func TestFixedAsserts(t *testing.T) {
	as := assert.New(t)
	as.Equal(10, arity.AssertFixed(10, 10))
	defer as.ExpectPanic(fmt.Sprintf(arity.ErrFixedArity, 9, 10))
	arity.AssertFixed(9, 10)
}

func TestMinimumAsserts(t *testing.T) {
	as := assert.New(t)
	as.Equal(5, arity.AssertMinimum(5, 5))
	defer as.ExpectPanic(fmt.Sprintf(arity.ErrMinimumArity, 10, 9))
	arity.AssertMinimum(10, 9)
}

func TestRangedAsserts(t *testing.T) {
	as := assert.New(t)
	as.Equal(5, arity.AssertRanged(3, 7, 5))
	defer as.ExpectPanic(fmt.Sprintf(arity.ErrRangedArity, 3, 7, 2))
	arity.AssertRanged(3, 7, 2)
}

func TestMakeChecker(t *testing.T) {
	as := assert.New(t)
	as.Nil(arity.MakeChecker())

	fn1 := arity.MakeChecker(1)
	as.Nil(fn1(1))
	as.Errorf(fn1(2), arity.ErrFixedArity, 1, 2)

	fn2 := arity.MakeChecker(2, arity.OrMore)
	as.Nil(fn2(5))
	as.Errorf(fn2(1), arity.ErrMinimumArity, 2, 1)

	fn3 := arity.MakeChecker(2, 7)
	as.Nil(fn3(4))
	as.Errorf(fn3(8), arity.ErrRangedArity, 2, 7, 8)

	defer as.ExpectPanic("too many arity check arguments")
	arity.MakeChecker(1, 2, 3)
}
