package main_test

import (
	"testing"

	main "github.com/kode4food/ale/cmd/ale"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestDocString(t *testing.T) {
	as := assert.New(t)

	ifStr, err := main.GetDocString("if")
	as.Contains("---", S(ifStr))
	as.Nil(err)

	s, err := main.GetDocString("no-way-this-exists")
	as.Empty(s)
	as.Errorf(err, main.ErrDocNotFound, "no-way-this-exists")
}
