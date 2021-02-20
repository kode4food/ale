package encoder_test

import (
	"testing"

	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestEncoderEquality(t *testing.T) {
	as := NewWrapped(t)
	e1 := getTestEncoder()
	e2 := getTestEncoder()
	as.True(e1.Equal(e1))
	as.False(e1.Equal(e2))
	as.False(e1.Equal(S("hello")))
}
