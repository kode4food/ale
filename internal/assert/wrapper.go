package assert

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/stretchr/testify/assert"
)

type (
	// Any is the friendly name for a generic interface
	Any interface{}

	// Wrapper wraps the testify assertions module in order to perform
	// checking and conversion that is system-specific
	Wrapper struct {
		*testing.T
		*assert.Assertions
	}
)

// Error messages
const (
	ErrInvalidTestExpression = "invalid test expression: %v"
	ErrProperErrorNotRaised  = "proper error not raised"
	ErrValueNotFound         = "value not found in object: %s"
)

// New instantiates a new Wrapper instance from the specified test
func New(t *testing.T) *Wrapper {
	return &Wrapper{
		T:          t,
		Assertions: assert.New(t),
	}
}

// String tests a Value for string equality
func (w *Wrapper) String(expect string, expr Any) {
	w.Helper()
	switch s := expr.(type) {
	case string:
		w.Assertions.Equal(expect, s)
	case data.Value:
		w.Assertions.Equal(expect, s.String())
	default:
		panic(fmt.Errorf(ErrInvalidTestExpression, expr))
	}
}

// Number tests a Value for numeric equality
func (w *Wrapper) Number(expect float64, expr Any) {
	w.Helper()
	switch n := expr.(type) {
	case float64:
		w.Assertions.Equal(expect, n)
	case int:
		w.Assertions.Equal(int64(expect), int64(n))
	case data.Number:
		w.Assertions.Equal(data.EqualTo, data.Float(expect).Cmp(n))
	default:
		panic(fmt.Errorf(ErrInvalidTestExpression, expr))
	}
}

// Equal tests a Value for some kind of equality. Performs checks to do so
func (w *Wrapper) Equal(expect Any, expr Any) {
	w.Helper()
	switch expect := expect.(type) {
	case data.String:
		w.String(string(expect), expr)
	case data.Number:
		num := expr.(data.Number)
		w.Assertions.Equal(data.EqualTo, expect.Cmp(num))
	case data.Value:
		if expr, ok := expr.(data.Value); ok {
			w.True(expect.Equal(expr))
		} else {
			w.String(expect.String(), expr)
		}
	default:
		w.Assertions.Equal(expect, expr)
	}
}

// True tests a Value for boolean true
func (w *Wrapper) True(expr Any) {
	w.Helper()
	if b, ok := expr.(data.Bool); ok {
		w.Assertions.True(bool(b))
		return
	}
	w.Assertions.True(expr.(bool))
}

// Truthy tests a Value for system-specific Truthy
func (w *Wrapper) Truthy(expr data.Value) {
	w.Helper()
	w.Assertions.True(data.Truthy(expr))
}

// False tests a Value for boolean false
func (w *Wrapper) False(expr Any) {
	w.Helper()
	if b, ok := expr.(data.Bool); ok {
		w.Assertions.False(bool(b))
		return
	}
	w.Assertions.False(expr.(bool))
}

// Falsey tests a Value for system-specific Falsey
func (w *Wrapper) Falsey(expr data.Value) {
	w.Helper()
	w.Assertions.False(data.Truthy(expr))
}

// Contains check if the expected string is in the provided Value
func (w *Wrapper) Contains(expect string, expr data.Value) {
	w.Helper()
	val := expr.String()
	w.Assertions.True(strings.Contains(val, expect))
}

// NotContains checks if the expected string is not in the provided Value
func (w *Wrapper) NotContains(expect string, expr data.Value) {
	w.Helper()
	val := expr.String()
	w.Assertions.False(strings.Contains(val, expect))
}

// Identical tests that two values are referentially identical
func (w *Wrapper) Identical(expect Any, expr Any) {
	w.Helper()
	w.Assertions.Equal(expect, expr)
}

// NotIdentical tests that two values are not referentially identical
func (w *Wrapper) NotIdentical(expect Any, expr Any) {
	w.Helper()
	w.Assertions.NotEqual(expect, expr)
}

// Compare tests if the Comparison of two Numbers is correct
func (w *Wrapper) Compare(c data.Comparison, l data.Number, r data.Number) {
	w.Helper()
	w.Assertions.Equal(c, l.Cmp(r))
}

// ExpectPanic is used with a defer to make sure an error was triggered
func (w *Wrapper) ExpectPanic(errStr string) {
	w.Helper()
	if rec := recover(); rec != nil {
		if re, ok := rec.(error); ok {
			recStr := re.Error()
			w.True(strings.HasPrefix(recStr, errStr))
			return
		}
	}
	panic(ErrProperErrorNotRaised)
}

// ExpectNoPanic is used with defer to make sure no error occurs
func (w *Wrapper) ExpectNoPanic() {
	w.Helper()
	rec := recover()
	w.Nil(rec)
}

// MustGet retrieves a Value from a Mapped or explodes
func (w *Wrapper) MustGet(m data.Mapped, k data.Value) data.Value {
	if v, ok := m.Get(k); ok {
		return v
	}
	panic(fmt.Errorf(ErrValueNotFound, k))
}
