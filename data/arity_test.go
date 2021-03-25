package data_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
)

func TestFixedAsserts(t *testing.T) {
	as := assert.New(t)
	as.Equal(10, data.AssertFixed(10, 10))
	defer as.ExpectPanic(fmt.Sprintf(data.ErrFixedArity, 9, 10))
	data.AssertFixed(9, 10)
}

func TestMinimumAsserts(t *testing.T) {
	as := assert.New(t)
	as.Equal(5, data.AssertMinimum(5, 5))
	defer as.ExpectPanic(fmt.Sprintf(data.ErrMinimumArity, 10, 9))
	data.AssertMinimum(10, 9)
}

func TestRangedAsserts(t *testing.T) {
	as := assert.New(t)
	as.Equal(5, data.AssertRanged(3, 7, 5))
	defer as.ExpectPanic(fmt.Sprintf(data.ErrRangedArity, 3, 7, 2))
	data.AssertRanged(3, 7, 2)
}

func TestMakeChecker(t *testing.T) {
	as := assert.New(t)
	as.Nil(data.MakeChecker())

	fn1 := data.MakeChecker(1)
	as.Nil(fn1(1))
	as.EqualError(fn1(2), fmt.Sprintf(data.ErrFixedArity, 1, 2))

	fn2 := data.MakeChecker(2, data.OrMore)
	as.Nil(fn2(5))
	as.EqualError(fn2(1), fmt.Sprintf(data.ErrMinimumArity, 2, 1))

	fn3 := data.MakeChecker(2, 7)
	as.Nil(fn3(4))
	as.EqualError(fn3(8), fmt.Sprintf(data.ErrRangedArity, 2, 7, 8))

	defer as.ExpectPanic("too many arity check arguments")
	data.MakeChecker(1, 2, 3)
}
