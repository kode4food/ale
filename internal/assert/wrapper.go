package assert

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/debug"
	"github.com/stretchr/testify/assert"
)

// Wrapper wraps testify assertions to perform checking and conversion that is
// system-specific
type Wrapper struct {
	*testing.T
	*assert.Assertions
}

const (
	// ErrInvalidTestExpression is raised when a check encounters a test value
	// that is not expected
	ErrInvalidTestExpression = "invalid test expression: %v"

	// ErrProperErrorNotRaised is raised when a panic is expected but not seen
	ErrProperErrorNotRaised = "proper error not raised, expected: %s"

	// ErrCannotMakeString is raised when a value can't be made into a string
	ErrCannotMakeString = "can't convert value to string"

	// ErrValueNotFound is raised when forced retrieval from an Object fails
	ErrValueNotFound = "value not found in object: %s"
)

// New instantiates a new Wrapper instance from the specified test
func New(t *testing.T) *Wrapper {
	return &Wrapper{
		T:          t,
		Assertions: assert.New(t),
	}
}

// String tests a Value for string equality
func (w *Wrapper) String(expect string, expr any) {
	w.Helper()
	switch s := expr.(type) {
	case string:
		w.Assertions.Equal(expect, s)
	case data.Local:
		w.Assertions.Equal(expect, string(s))
	case ale.Value:
		w.Assertions.Equal(expect, data.ToString(s))
	default:
		panic(debug.ProgrammerErrorf(ErrInvalidTestExpression, expr))
	}
}

// Number tests a Value for numeric equality
func (w *Wrapper) Number(expect float64, expr any) {
	w.Helper()
	switch n := expr.(type) {
	case float64:
		w.Assertions.Equal(expect, n)
	case int:
		w.Assertions.Equal(int64(expect), int64(n))
	case data.Number:
		w.Compare(data.EqualTo, data.Float(expect), n)
	default:
		panic(debug.ProgrammerErrorf(ErrInvalidTestExpression, expr))
	}
}

// Equal tests a Value for some kind of equality
func (w *Wrapper) Equal(expect any, expr any) {
	w.Helper()
	switch expect := expect.(type) {
	case data.String:
		w.String(string(expect), expr)
	case data.Number:
		num := expr.(data.Number)
		w.Assertions.Equal(data.EqualTo, expect.Cmp(num))
	case ale.Value:
		if expr, ok := expr.(ale.Value); ok {
			w.True(expect.Equal(expr))
			return
		}
		w.String(data.ToString(expect), expr)
	default:
		w.Assertions.Equal(expect, expr)
	}
}

// True tests a Value for boolean true
func (w *Wrapper) True(expr any) {
	w.Helper()
	if b, ok := expr.(data.Bool); ok {
		w.Assertions.True(bool(b))
		return
	}
	w.Assertions.True(expr.(bool))
}

// False tests a Value for boolean false
func (w *Wrapper) False(expr any) {
	w.Helper()
	if b, ok := expr.(data.Bool); ok {
		w.Assertions.False(bool(b))
		return
	}
	w.Assertions.False(expr.(bool))
}

// Contains check if the expected string is in the provided Value
func (w *Wrapper) Contains(expect string, expr ale.Value) {
	w.Helper()
	val := data.ToString(expr)
	w.Assertions.True(strings.Contains(val, expect))
}

// NotContains checks if the expected string is not in the provided Value
func (w *Wrapper) NotContains(expect string, expr ale.Value) {
	w.Helper()
	val := data.ToString(expr)
	w.Assertions.False(strings.Contains(val, expect))
}

// Identical tests that two values are referentially identical
func (w *Wrapper) Identical(expect any, expr any) {
	w.Helper()
	p1 := fmt.Sprintf("%p", expect)
	p2 := fmt.Sprintf("%p", expr)
	w.Assertions.True(p1 == p2)
}

// NotIdentical tests that two values are not referentially identical
func (w *Wrapper) NotIdentical(expect any, expr any) {
	w.Helper()
	p1 := fmt.Sprintf("%p", expect)
	p2 := fmt.Sprintf("%p", expr)
	w.Assertions.False(p1 == p2)
}

// Compare tests if the Comparison of two Numbers is correct
func (w *Wrapper) Compare(c data.Comparison, l data.Number, r data.Number) {
	w.Helper()
	w.Assertions.Equal(c, l.Cmp(r))
}

func (w *Wrapper) ExpectError(expected, err any) {
	errStr := w.mustMakeString(err)
	pfx := w.mustMakeString(expected)
	hasPfx := strings.HasPrefix(errStr, pfx)
	w.True(hasPfx)
	if rec, ok := err.(error); ok && !hasPfx {
		w.EqualError(rec, pfx)
	}
}

// ExpectProgrammerError is used with defer to make sure a programmer error
// was triggered
func (w *Wrapper) ExpectProgrammerError(errStr string) {
	w.Helper()
	if rec := recover(); rec != nil {
		if recStr, ok := rec.(string); ok {
			w.Equal(errStr, recStr)
			return
		}
	}
	w.Fail(ErrProperErrorNotRaised)
}

// MustGet retrieves a Value from a Mapped or explodes
func (w *Wrapper) MustGet(m data.Mapped, k ale.Value) ale.Value {
	if v, ok := m.Get(k); ok {
		return v
	}
	panic(fmt.Errorf(ErrValueNotFound, k))
}

func (w *Wrapper) mustMakeString(val any) string {
	res, err := w.makeString(val)
	if err != nil {
		panic(debug.ProgrammerErrorf("%w", err))
	}
	return res
}

func (w *Wrapper) makeString(val any) (string, error) {
	switch val := val.(type) {
	case string:
		return val, nil
	case error:
		return val.Error(), nil
	case fmt.Stringer:
		return val.String(), nil
	default:
		return "", errors.New(ErrCannotMakeString)
	}
}
