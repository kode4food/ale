package internal_test

import (
	"testing"

	"github.com/kode4food/ale/cmd/ale/internal"
	"github.com/kode4food/ale/cmd/ale/internal/console"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestEmptyContainers(t *testing.T) {
	as := assert.New(t)

	as.Equal("()", internal.PrettyPrintAt((*data.List)(nil), 0))
	as.Equal("[]", internal.PrettyPrintAt(data.EmptyVector, 0))
	as.Equal("{}", internal.PrettyPrintAt(data.EmptyObject, 0))
	as.Equal("{}", internal.PrettyPrintAt((*data.Object)(nil), 0))
}

func TestInlineContainers(t *testing.T) {
	as := assert.New(t)

	originalGetScreenWidth := console.GetScreenWidth
	console.GetScreenWidth = func() int { return 120 }
	defer func() { console.GetScreenWidth = originalGetScreenWidth }()

	as.Equal("(1 2 3)", internal.PrettyPrintAt(L(I(1), I(2), I(3)), 0))
	as.Equal(
		"[1 2 3]", internal.PrettyPrintAt(data.NewVector(I(1), I(2), I(3)), 0),
	)
}

func TestMultilineContainers(t *testing.T) {
	as := assert.New(t)

	originalGetScreenWidth := console.GetScreenWidth
	console.GetScreenWidth = func() int { return 1 }
	defer func() { console.GetScreenWidth = originalGetScreenWidth }()

	as.Equal(
		"(\n  1\n  2\n  3\n)", internal.PrettyPrintAt(L(I(1), I(2), I(3)), 0),
	)
	as.Equal(
		"[\n  1\n  2\n  3\n]",
		internal.PrettyPrintAt(data.NewVector(I(1), I(2), I(3)), 0),
	)
}

func TestObjects(t *testing.T) {
	as := assert.New(t)

	obj1 := data.NewObject(C(K("key"), S("value")))
	as.Equal(`{:key "value"}`, internal.PrettyPrintAt(obj1, 0))

	obj2 := data.NewObject(C(K("name"), S("Thom")), C(K("year"), I(2025)))
	as.Equal(
		"{\n  :name \"Thom\"\n  :year 2025\n}", internal.PrettyPrintAt(obj2, 0),
	)

	obj3 := data.NewObject(C(K("short"), I(1)), C(K("longer"), I(2)))
	as.Equal(
		"{\n  :longer 2\n  :short  1\n}", internal.PrettyPrintAt(obj3, 0),
	)
}

func TestCons(t *testing.T) {
	as := assert.New(t)

	as.Equal("(1 . 2)", internal.PrettyPrintAt(C(I(1), I(2)), 0))

	obj1 := data.NewObject(C(K("age"), I(63)), C(K("name"), S("Thom")))
	obj2 := data.NewObject(
		C(K("sound"), S("Hee-haw")),
		C(K("animal"), S("Donkey")),
	)
	expected := "({\n  :age  63\n  :name \"Thom\"\n} " +
		". {\n  :animal \"Donkey\"\n  :sound  \"Hee-haw\"\n})"
	as.Equal(expected, internal.PrettyPrintAt(C(obj1, obj2), 0))
}

func TestNestedStructures(t *testing.T) {
	as := assert.New(t)

	originalGetScreenWidth := console.GetScreenWidth
	console.GetScreenWidth = func() int { return 120 }
	defer func() { console.GetScreenWidth = originalGetScreenWidth }()

	as.Equal(
		"((1 2) (3 4))",
		internal.PrettyPrintAt(L(L(I(1), I(2)), L(I(3), I(4))), 0),
	)

	listWithObj := L(
		I(1), data.NewObject(C(K("key"), S("value")), C(K("other"), S("data"))),
	)
	expected := "(\n  1 {\n    :key   \"value\"\n    :other \"data\"\n  }\n)"
	as.Equal(expected, internal.PrettyPrintAt(listWithObj, 0))
}

func TestObjectWithComplexKeys(t *testing.T) {
	as := assert.New(t)

	originalGetScreenWidth := console.GetScreenWidth
	console.GetScreenWidth = func() int { return 120 }
	defer func() { console.GetScreenWidth = originalGetScreenWidth }()

	nestedKey := L(I(1), I(2), I(3))
	obj := data.NewObject(
		C(nestedKey, S("value")),
		C(K("simple"), S("simple-value")),
	)
	expected := "{\n  (1 2 3) \"value\"\n  :simple \"simple-value\"\n}"
	as.Equal(expected, internal.PrettyPrintAt(obj, 0))
}

func TestWithOffset(t *testing.T) {
	as := assert.New(t)

	obj := data.NewObject(C(K("key"), S("value")), C(K("other"), S("data")))

	expected := "{\n  :key   \"value\"\n  :other \"data\"\n}"
	as.Equal(expected, internal.PrettyPrintAt(obj, 0))

	expected = "{\n      :key   \"value\"\n      :other \"data\"\n    }"
	as.Equal(expected, internal.PrettyPrintAt(obj, 4))
}

func TestSequenceWrapping(t *testing.T) {
	as := assert.New(t)

	originalGetScreenWidth := console.GetScreenWidth
	console.GetScreenWidth = func() int { return 20 }
	defer func() { console.GetScreenWidth = originalGetScreenWidth }()

	longList := L(S("elem1"), S("elem2"), S("elem3"))
	expected := "(\n  \"elem1\" \"elem2\"\n  \"elem3\"\n)"
	as.Equal(expected, internal.PrettyPrintAt(longList, 0))

	longVector := data.NewVector(S("elem1"), S("elem2"), S("elem3"))
	expected = "[\n  \"elem1\" \"elem2\"\n  \"elem3\"\n]"
	as.Equal(expected, internal.PrettyPrintAt(longVector, 0))
}

func TestElementPacking(t *testing.T) {
	as := assert.New(t)

	originalGetScreenWidth := console.GetScreenWidth
	console.GetScreenWidth = func() int { return 8 }
	defer func() { console.GetScreenWidth = originalGetScreenWidth }()

	manyElements := L(I(1), I(2), I(3), I(4), I(5))
	expected := "(\n  1 2 3\n  4 5\n)"
	as.Equal(expected, internal.PrettyPrintAt(manyElements, 0))
}

func TestEdgeCases(t *testing.T) {
	as := assert.New(t)

	obj := data.NewObject(C(K(""), S("")))
	as.Equal(`{: ""}`, internal.PrettyPrintAt(obj, 0))

	originalGetScreenWidth := console.GetScreenWidth
	console.GetScreenWidth = func() int { return 120 }
	defer func() { console.GetScreenWidth = originalGetScreenWidth }()

	deep := L(L(L(I(1))))
	as.Equal("(((1)))", internal.PrettyPrintAt(deep, 0))

	cons := C(L(I(1), I(2)), data.NewVector(I(3), I(4)))
	as.Equal("((1 2) . [3 4])", internal.PrettyPrintAt(cons, 0))

	mixedObj := data.NewObject(
		C(K("list"), L(I(1), I(2))),
		C(K("vector"), data.NewVector(I(3), I(4))),
		C(K("string"), S("hello")),
	)
	expected := "{\n  :list   (1 2)\n  :string \"hello\"\n  :vector [3 4]\n}"
	as.Equal(expected, internal.PrettyPrintAt(mixedObj, 0))
}
