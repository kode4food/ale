package builtin_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/core/builtin"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/runtime"
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
	as.MustEvalTo(`(length {:name "Ale" :age 45})`, F(2))
	as.MustEvalTo(`(length (object :name "Ale" :age 45))`, F(2))
	as.MustEvalTo(`(object? {:name "Ale" :age 45})`, data.True)
	as.MustEvalTo(`(object? (object :name "Ale" :age 45))`, data.True)
	as.MustEvalTo(`(object? '(:name "Ale" :age 45))`, data.False)
	as.MustEvalTo(`(object? [:name "Ale" :age 45])`, data.False)
	as.MustEvalTo(`(!object? '(:name "Ale" :age 45))`, data.True)
	as.MustEvalTo(`(!object? [:name "Ale" :age 45])`, data.True)
	as.MustEvalTo(`(:name {:name "Ale" :age 45})`, S("Ale"))

	as.MustEvalTo(`
		(:name (apply object (concat '(:name "Ale") '(:age 45))))
	`, S("Ale"))

	as.MustEvalTo(`
		(define x {:name "bob" :age 45})
		(x :name)
	`, S("bob"))

	as.PanicWith(`(object :too "few" :args)`, data.ErrMapNotPaired)

	as.PanicWith(`
		(apply object (concat '(:name "Ale") '(:age)))
	`, data.ErrMapNotPaired)
}

func TestObjectAssoc(t *testing.T) {
	as := assert.New(t)

	as.MustEvalTo(`
		(define o1 {:first "first" :second "second"})
		(define o2 (assoc o1 (:first . "first-replaced")))
		(define o3 (assoc o1 (:first . "also-replaced")))
		(define o4 (dissoc o1 :first))
		(define o5 (dissoc {} :first))
		[(:first o1) (:second o1)
         (:first o2) (:second o2)
         (:first o3) (:second o3)
         (:first o4) (:second o4)
		 (:first o5) (:second o5)]
	`,
		V(
			S("first"), S("second"),
			S("first-replaced"), S("second"),
			S("also-replaced"), S("second"),
			data.Null, S("second"),
			data.Null, data.Null,
		),
	)

	as.PanicWith(`
		(assoc {} :not-a-pair)
	`, fmt.Errorf(runtime.ErrUnexpectedType, "keyword", "pair"))

}

func TestMappedEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`(mapped? {:name "Ale" :age 45})`, data.True)
	as.MustEvalTo(`(mapped? (object :name "Ale" :age 45))`, data.True)
	as.MustEvalTo(`(mapped? '(:name "Ale" :age 45))`, data.False)
	as.MustEvalTo(`(mapped? [:name "Ale" :age 45])`, data.False)
	as.MustEvalTo(`(!mapped? '(:name "Ale" :age 45))`, data.True)
	as.MustEvalTo(`(!mapped? '(:name "Ale" :age 45))`, data.True)
	as.MustEvalTo(`(!mapped? [:name "Ale" :age 45])`, data.True)
}
