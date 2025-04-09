package assert_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/sync"
	"github.com/kode4food/ale/pkg/data"
)

func TestTheStringTests(t *testing.T) {
	as := assert.New(t)

	as.String("hello", "hello")
	as.String("hello", S("hello"))
	as.String(":hello", K("hello"))
	as.String("hello", LS("hello"))

	defer as.ExpectPanic(fmt.Errorf(assert.ErrInvalidTestExpression, "10"))
	as.String("10", 10)
}

func TestTheFloatTests(t *testing.T) {
	as := assert.New(t)

	as.Number(10.5, F(10.5))
	as.Number(10, F(10))
	as.Number(10, 10.0)
	as.Number(10, 10)

	defer as.ExpectPanic(fmt.Errorf(assert.ErrInvalidTestExpression, "10"))
	as.Number(10, "10")
}

func TestEquality(t *testing.T) {
	as := assert.New(t)

	as.Equal(S("hello"), S("hello"))
	as.Equal(S("hello"), "hello")
	as.NotEqual(S("hello"), S("goodbye"))

	as.Equal(I(37), I(37))
	as.NotEqual(I(37), F(42.0))

	as.Equal(L(I(1), I(2), I(3)), "(1 2 3)")
}

func TestBools(t *testing.T) {
	as := assert.New(t)

	as.True(data.True)
	as.True(B(true))
	as.True(true)

	as.False(data.False)
	as.False(B(false))
	as.False(false)
}

func TestContains(t *testing.T) {
	as := assert.New(t)

	as.Contains("New York", S("A New York State of Mind"))
	as.NotContains("Boston", S("A New York State of Mind"))

	as.Contains(":type promise", sync.NewPromise(nil))
}

func TestIdentical(t *testing.T) {
	as := assert.New(t)

	l1 := L(I(1), I(2), I(3))
	l2 := L(I(1), I(2), I(3))
	as.Identical(l1, l1)
	as.Equal(l1, l1)
	as.NotIdentical(l1, l2)
	as.Equal(l1, l2)
}

func TestMustGetExplosion(t *testing.T) {
	as := assert.New(t)
	// Can handle errors in multiple forms
	err := S(fmt.Sprintf(assert.ErrValueNotFound, K("hello")))
	defer as.ExpectPanic(err)
	as.MustGet(O(), K("hello"))
}

func TestMustGetNonExplosion(t *testing.T) {
	as := assert.New(t)
	defer as.ExpectNoPanic()
	as.MustGet(O(C(K("hello"), S("world"))), K("hello"))
}

func TestExpectProgrammerError(t *testing.T) {
	as := assert.New(t)
	defer as.ExpectProgrammerError("just a string")
	panic("just a string")
}
