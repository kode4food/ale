package helpers_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestHelpers(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(
		`'(1 1/2 hello
			{ale/hello "string"}
			{:kwd ["a" "vector"]}
			#t (1 . 2.3))
		`,
		L(
			I(1), R(1, 2), LS("hello"),
			O(C(QS("ale", "hello"), S("string"))),
			O(C(K("kwd"), V(S("a"), S("vector")))),
			B(true), C(I(1), F(2.3)),
		),
	)
}
