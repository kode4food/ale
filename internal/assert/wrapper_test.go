package assert_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestTheStringTests(t *testing.T) {
	as := assert.New(t)

	as.String("hello", "hello")
	as.String("hello", S("hello"))
	as.String(":hello", K("hello"))

	defer as.ExpectPanic(fmt.Sprintf(assert.InvalidTestExpression, "10"))
	as.String("10", 10)
}

func TestTheFloatTests(t *testing.T) {
	as := assert.New(t)

	as.Number(10.5, F(10.5))
	as.Number(10, F(10))

	defer as.ExpectPanic(fmt.Sprintf(assert.InvalidTestExpression, "10"))
	as.Number(10, "10")
}

func TestTheNonExplosions(t *testing.T) {
	as := assert.New(t)
	defer func() {
		if rec := recover(); rec != nil {
			as.String("proper error not raised", rec)
		}
	}()
	defer as.ExpectPanic("will not happen")
}
