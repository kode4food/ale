package internal

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/internal/lang/parse"
	"github.com/kode4food/ale/pkg/env"
)

func panicWithRec(t *testing.T, recoverable bool, msg string, fn func()) {
	as := assert.New(t)
	defer func() {
		if err := toError(recover()); err != nil {
			as.NotNil(err)
			as.Equal(recoverable, isRecoverable(err))
			as.ErrorContains(err, msg)
		}
	}()
	fn()
}

func panicWith(t *testing.T, msg string, fn func()) {
	panicWithRec(t, false, msg, fn)
}

func panicWithRecoverable(t *testing.T, msg string, fn func()) {
	panicWithRec(t, true, msg, fn)
}

func TestReplEvalBuffer(t *testing.T) {
	as := assert.New(t)

	r := NewREPL()
	r.buf.WriteString("(+ 1 2")
	as.False(r.evalBuffer())

	r.buf.WriteString(" 3)")
	as.True(r.evalBuffer())
}

func TestEvalBuffer(t *testing.T) {
	as := assert.New(t)

	as.Nil(evalBuffer([]byte(`"hello world"`)))
	as.EqualError(
		evalBuffer([]byte(`(unknown)`)),
		fmt.Sprintf(env.ErrNameNotDeclared, "unknown"),
	)

	panicWith(t, "boom", func() {
		mustEvalBuffer([]byte(`(raise "boom")`))
	})

	panicWithRecoverable(t, parse.ErrListNotClosed, func() {
		mustEvalBuffer([]byte(`(no-close `))
	})
}
