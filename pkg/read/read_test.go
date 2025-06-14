package read_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/read"
)

func TestFromString(t *testing.T) {
	as := assert.New(t)
	ns := assert.GetTestNamespace()
	tr := read.FromString(ns, "99")
	as.NotNil(tr)
	as.Equal(I(99), tr.Car())
}
