package data_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestLastOfSequence(t *testing.T) {
	as := assert.New(t)

	v, ok := data.Last(data.EmptyList)
	as.Nil(v)
	as.False(ok)

	v, ok = data.Last(L(S("this"), S("is"), S("last")))
	as.String("last", v)
	as.True(ok)

	v, ok = data.Last(V(S("this"), S("is"), S("last")))
	as.String("last", v)
	as.True(ok)
}
