package builtin_test

import (
	"fmt"
	"testing"

	"gitlab.com/kode4food/ale/bootstrap/builtin"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestAssoc(t *testing.T) {
	as := assert.New(t)

	a1 := builtin.Assoc(K("hello"), S("foo"))
	m1 := a1.(data.Mapped)
	v1, ok := m1.Get(K("hello"))
	as.True(ok)
	as.String("foo", v1)

	as.True(builtin.IsAssoc(a1))
	as.False(builtin.IsAssoc(I(99)))

	as.True(builtin.IsMapped(a1))
	as.False(builtin.IsMapped(I(99)))
}

func TestAssocEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(len {:name "Ale", :age 45})`, F(2))
	as.EvalTo(`(len (assoc :name "Ale", :age 45))`, F(2))
	as.EvalTo(`(assoc? {:name "Ale" :age 45})`, data.True)
	as.EvalTo(`(assoc? (assoc :name "Ale" :age 45))`, data.True)
	as.EvalTo(`(assoc? '(:name "Ale" :age 45))`, data.False)
	as.EvalTo(`(assoc? [:name "Ale" :age 45])`, data.False)
	as.EvalTo(`(!assoc? '(:name "Ale" :age 45))`, data.True)
	as.EvalTo(`(!assoc? [:name "Ale" :age 45])`, data.True)
	as.EvalTo(`(:name {:name "Ale" :age 45})`, S("Ale"))

	as.EvalTo(`
		(:name (apply assoc (concat '(:name "Ale") '(:age 45))))
	`, S("Ale"))

	as.EvalTo(`
		(def x {:name "bob" :age 45})
		(x :name)
	`, S("bob"))

	as.PanicWith(`(assoc :too "few" :args)`, fmt.Errorf(data.ExpectedPair))

	as.PanicWith(`
		(apply assoc (concat '(:name "Ale") '(:age)))
	`, fmt.Errorf(data.ExpectedPair))
}

func TestMappedEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(mapped? {:name "Ale" :age 45})`, data.True)
	as.EvalTo(`(mapped? (assoc :name "Ale" :age 45))`, data.True)
	as.EvalTo(`(mapped? '(:name "Ale" :age 45))`, data.False)
	as.EvalTo(`(mapped? [:name "Ale" :age 45])`, data.False)
	as.EvalTo(`(!mapped? '(:name "Ale" :age 45))`, data.True)
	as.EvalTo(`(!mapped? '(:name "Ale" :age 45))`, data.True)
	as.EvalTo(`(!mapped? [:name "Ale" :age 45])`, data.True)
}
