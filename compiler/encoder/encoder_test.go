package encoder_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestEncoderEquality(t *testing.T) {
	as := assert.New(t)
	e1 := assert.GetTestEncoder()
	e2 := assert.GetTestEncoder()
	as.True(e1.Equal(e1))
	as.False(e1.Equal(e2))
	as.False(e1.Equal(S("hello")))
}
