package special_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
)

func TestAsm(t *testing.T) {
	as := assert.New(t)
	as.Eval(`(asm* (emit :return))`)
}
