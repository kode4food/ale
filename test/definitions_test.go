package test

import (
	"testing"

	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestFunction(t *testing.T) {
	testCode(t, `
		(defn say-hello
		  []
		  "Hello, World!")
		(say-hello)
	`, S("Hello, World!"))

	testCode(t, `
		(defn identity [value] value)
		(identity "foo")
	`, S("foo"))
}

func TestBadFunction(t *testing.T) {
	symErr := intfErr("api.Integer", "api.LocalSymbol", "LocalSymbol")
	listErr := typeErr("api.Integer", "*api.List")
	vecErr := typeErr("api.Integer", "api.Vector")

	testBadCode(t, `(defn blah [name 99 bad] (name))`, symErr)
	testBadCode(t, `(defn blah 99 (name))`, listErr)
	testBadCode(t, `(defn 99 [x y] (+ x y))`, symErr)
	testBadCode(t, `(defn blah (99 "hello"))`, vecErr)
	testBadCode(t, `(defn blah ([x] "hello") 99)`, listErr)
}
