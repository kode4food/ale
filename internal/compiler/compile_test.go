package compiler_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/compiler"
	"github.com/kode4food/ale/pkg/data"
)

func TestIsEvaluable(t *testing.T) {
	as := assert.New(t)

	as.True(compiler.IsEvaluable(L(LS("some-sym"))))
	as.False(compiler.IsEvaluable(S("some-string")))
	as.True(compiler.IsEvaluable(LS("some-sym")))
	as.False(compiler.IsEvaluable(C(K("keyword"), S("some-value"))))
	as.False(compiler.IsEvaluable(data.True))
}
