package data_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestNames(t *testing.T) {
	as := assert.New(t)

	n := LS("hello")
	as.Equal(LS("hello"), n)
}
