package data_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/pkg/data"
)

func TestMakeChecker(t *testing.T) {
	as := assert.New(t)
	fn1, err := data.MakeArityChecker()
	as.NoError(err)
	as.NotNil(fn1)
	as.Nil(fn1(-1))
	as.Nil(fn1(1000))

	fn2, err := data.MakeArityChecker(1)
	as.NoError(err)
	as.Nil(fn2(1))
	as.EqualError(fn2(2), fmt.Sprintf(data.ErrFixedArity, 1, 2))

	fn3, err := data.MakeArityChecker(2, data.OrMore)
	as.NoError(err)
	as.Nil(fn3(5))
	as.EqualError(fn3(1), fmt.Sprintf(data.ErrMinimumArity, 2, 1))

	fn4, err := data.MakeArityChecker(2, 7)
	as.NoError(err)
	as.Nil(fn4(4))
	as.EqualError(fn4(8), fmt.Sprintf(data.ErrRangedArity, 2, 7, 8))

	_, err = data.MakeArityChecker(1, 2, 3)
	as.EqualError(err, data.ErrTooManyArguments)
}
