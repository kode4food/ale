package test

import (
	"fmt"
	"testing"

	"gitlab.com/kode4food/ale/compiler/special"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestDefinitions(t *testing.T) {
	testCode(t, `
		(def foo "bar")
		foo
	`, S("bar"))

	testCode(t, `
		(def return-local (fn []
			(let [foo "local"] foo)))
		(return-local)
	`, S("local"))
}

func TestLetBindings(t *testing.T) {
	testBadCode(t, `
		(let 99 "hello")
	`, typeErr("api.Integer", "api.Vector"))

	testBadCode(t, `
		(let [a blah b] "hello")
	`, fmt.Errorf(special.UnpairedBindings))
}
