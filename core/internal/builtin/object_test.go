package builtin_test

import (
	"errors"
	"testing"

	"github.com/kode4food/ale/core/internal/builtin"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestObject(t *testing.T) {
	as := assert.New(t)

	a1 := builtin.Object.Call(K("hello"), S("foo"))
	m1 := a1.(data.Mapped)
	v1, ok := m1.Get(K("hello"))
	as.True(ok)
	as.String("foo", v1)

	as.True(getPredicate(builtin.ObjectKey).Call(a1))
	as.False(getPredicate(builtin.ObjectKey).Call(I(99)))

	as.True(getPredicate(builtin.MappedKey).Call(a1))
	as.False(getPredicate(builtin.MappedKey).Call(I(99)))
}

func TestObjectEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(length {:name "Ale" :age 45})`, F(2))
	as.EvalTo(`(length (object :name "Ale" :age 45))`, F(2))
	as.EvalTo(`(object? {:name "Ale" :age 45})`, data.True)
	as.EvalTo(`(object? (object :name "Ale" :age 45))`, data.True)
	as.EvalTo(`(object? '(:name "Ale" :age 45))`, data.False)
	as.EvalTo(`(object? [:name "Ale" :age 45])`, data.False)
	as.EvalTo(`(!object? '(:name "Ale" :age 45))`, data.True)
	as.EvalTo(`(!object? [:name "Ale" :age 45])`, data.True)
	as.EvalTo(`(:name {:name "Ale" :age 45})`, S("Ale"))

	as.EvalTo(`
		(:name (apply object (concat '(:name "Ale") '(:age 45))))
	`, S("Ale"))

	as.EvalTo(`
		(define x {:name "bob" :age 45})
		(x :name)
	`, S("bob"))

	as.PanicWith(`(object :too "few" :args)`, errors.New(data.ErrMapNotPaired))

	as.PanicWith(`
		(apply object (concat '(:name "Ale") '(:age)))
	`, errors.New(data.ErrMapNotPaired))
}

func TestMappedEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(mapped? {:name "Ale" :age 45})`, data.True)
	as.EvalTo(`(mapped? (object :name "Ale" :age 45))`, data.True)
	as.EvalTo(`(mapped? '(:name "Ale" :age 45))`, data.False)
	as.EvalTo(`(mapped? [:name "Ale" :age 45])`, data.False)
	as.EvalTo(`(!mapped? '(:name "Ale" :age 45))`, data.True)
	as.EvalTo(`(!mapped? '(:name "Ale" :age 45))`, data.True)
	as.EvalTo(`(!mapped? [:name "Ale" :age 45])`, data.True)
}
