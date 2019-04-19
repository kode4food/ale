package test

import (
	"fmt"
	"testing"

	"gitlab.com/kode4food/ale/data"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestAssoc(t *testing.T) {
	testCode(t, `(len {:name "Ale", :age 45})`, F(2))
	testCode(t, `(len (assoc :name "Ale", :age 45))`, F(2))
	testCode(t, `(assoc? {:name "Ale" :age 45})`, data.True)
	testCode(t, `(assoc? (assoc :name "Ale" :age 45))`, data.True)
	testCode(t, `(assoc? '(:name "Ale" :age 45))`, data.False)
	testCode(t, `(assoc? [:name "Ale" :age 45])`, data.False)
	testCode(t, `(!assoc? '(:name "Ale" :age 45))`, data.True)
	testCode(t, `(!assoc? [:name "Ale" :age 45])`, data.True)
	testCode(t, `(:name {:name "Ale" :age 45})`, S("Ale"))

	testCode(t, `
		(:name (apply assoc (concat '(:name "Ale") '(:age 45))))
	`, S("Ale"))

	testCode(t, `
		(def x {:name "bob" :age 45})
		(x :name)
	`, S("bob"))

	testBadCode(t, `(assoc :too "few" :args)`, fmt.Errorf(data.ExpectedPair))

	testBadCode(t, `
		(apply assoc (concat '(:name "Ale") '(:age)))
	`, fmt.Errorf(data.ExpectedPair))
}

func TestMapped(t *testing.T) {
	testCode(t, `(mapped? {:name "Ale" :age 45})`, data.True)
	testCode(t, `(mapped? (assoc :name "Ale" :age 45))`, data.True)
	testCode(t, `(mapped? '(:name "Ale" :age 45))`, data.False)
	testCode(t, `(mapped? [:name "Ale" :age 45])`, data.False)
	testCode(t, `(!mapped? '(:name "Ale" :age 45))`, data.True)
	testCode(t, `(!mapped? '(:name "Ale" :age 45))`, data.True)
	testCode(t, `(!mapped? [:name "Ale" :age 45])`, data.True)
}
