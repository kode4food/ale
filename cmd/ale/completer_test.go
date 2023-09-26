package main_test

import (
	"testing"

	main "github.com/kode4food/ale/cmd/ale"
	"github.com/kode4food/ale/internal/assert"
)

func TestREPLDoBasic(t *testing.T) {
	as := assert.New(t)

	r := main.NewREPL()
	res, idx := r.Do([]rune("(seq->v'(1 2 3 4)"), 7)
	as.Equal(0, idx)
	as.Equal(1, len(res))
	as.Equal("ector ", string(res[0]))

	res, idx = r.Do([]rune("(seq->v '(1 2 3 4)"), 7)
	as.Equal(0, idx)
	as.Equal(1, len(res))
	as.Equal("ector", string(res[0]))

	res, idx = r.Do([]rune("(seq->"), 6)
	as.Equal(0, idx)
	as.Equal(3, len(res))
	as.Equal("list ", string(res[0]))
	as.Equal("object ", string(res[1]))
	as.Equal("vector ", string(res[2]))
}

func TestREPLDoNamespaced(t *testing.T) {
	as := assert.New(t)

	r := main.NewREPL()
	res, idx := r.Do([]rune("(al'(1 2 3 4)"), 3)
	as.Equal(0, idx)
	as.Equal(1, len(res))
	as.Equal("e/", string(res[0]))

	res, idx = r.Do([]rune("(ale/seq->"), 10)
	as.Equal(0, idx)
	as.Equal(3, len(res))
	as.Equal("list ", string(res[0]))
	as.Equal("object ", string(res[1]))
	as.Equal("vector ", string(res[2]))
}

func TestREPLDoEmptyResults(t *testing.T) {
	as := assert.New(t)

	r := main.NewREPL()
	res, idx := r.Do([]rune("(37'(1 2 3 4)"), 3)
	as.Equal(0, idx)
	as.Equal(0, len(res))

	res, idx = r.Do([]rune("(a/l/e'(1 2 3 4)"), 6)
	as.Equal(0, idx)
	as.Equal(0, len(res))
}
