package data_test

import (
	"testing"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestEmptyDeque(t *testing.T) {
	as := assert.New(t)

	d1 := data.EmptyDeque
	as.Equal(0, d1.Count())
	as.True(d1.IsEmpty())

	d1 = data.NewDeque(S("you"))
	d2 := d1.Prepend(S("oh")).(*data.Deque)
	d3 := d2.Append(S("pretty"), S("things")).(*data.Deque)
	as.String(`("oh" "you" "pretty" "things")`, d3)
	as.String("oh", d3.First())
	as.String(`("you" "pretty" "things")`, d3.Rest())
	as.String(`("things" "pretty" "you" "oh")`, d3.Reverse())
}
