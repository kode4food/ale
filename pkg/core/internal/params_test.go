package internal_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/pkg/core/internal"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/data"
)

func TestReachability(t *testing.T) {
	as := assert.New(t)

	as.PanicWith(`
		(lambda
			[(x y) "hello"]
			[(z) "there"]
			[(a b) "error"])
	`, fmt.Errorf(internal.ErrUnreachableCase, "(a b)"))

	as.PanicWith(`
		(lambda
			[(x y . z) "hello"]
			[(x y) "there"]
			[(a b) "error"])
	`, fmt.Errorf(internal.ErrUnreachableCase, "(x y)"))

	as.PanicWith(`
		(define-lambda test
			[(a b . c) #t]
		    [(a b c . d) #f])
	`, fmt.Errorf(internal.ErrUnreachableCase, "(a b c . d)"))

	as.EvalTo(`
		(define-lambda test
			[(a b c . d) #t]
		    [(a b . c) #f])
		[(test 1 2) (test 1 2 3)]
	`, V(data.False, data.True))

	as.EvalTo(`
		(define-lambda test
			[(a) a]
		    [(a b c d e f g h i) (+ a b c d e f g h i)])
		[(test 1 2 3 4 5 6 7 8 9) (test 1)]
	`, V(I(45), I(1)))
}

func TestUnmatchedCase(t *testing.T) {
	as := assert.New(t)

	as.PanicWith(`
		(define-lambda test
			[(a) [a]]
			[(a b) [a b]]
			[(a b c) [a b c]]
			[(a b c d e) [e]]
			[(a b c d e f g . h) [a]])
		(test)
	`, fmt.Errorf(internal.ErrUnmatchedCase, 0, "1-3, 5, 7 or more"))

	as.PanicWith(`
		(define-lambda test
			[(a) [a]]
			[(a b) [a b]]
			[(a b c) [a b c]]
			[(a b c d e) [e]]
			[(a b c d e g . h) [a]])
		(test)
	`, fmt.Errorf(internal.ErrUnmatchedCase, 0, "1-3, 5 or more"))

	as.PanicWith(`
		(define-lambda test
			[() #t]
			[(a) [a]]
			[(a b c) [a b c]]
			[(a b c d e) [e]])
		(test 1 2)
	`, fmt.Errorf(internal.ErrUnmatchedCase, 2, "0-1, 3, 5"))
}
