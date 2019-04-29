package assert

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/kode4food/ale/data"
)

type (
	// Any is the friendly name for a generic interface
	Any interface{}

	// Wrapper wraps the testify assertions module in order to perform
	// checking and conversion that is system-specific
	Wrapper struct {
		*assert.Assertions
	}
)

// Error messages
const (
	InvalidTestExpression = "invalid test expression: %v"
)

// New instantiates a new Wrapper instance from the specified test
func New(t *testing.T) *Wrapper {
	return &Wrapper{
		Assertions: assert.New(t),
	}
}

// String tests a Value for string equality
func (w *Wrapper) String(expect string, expr Any) {
	switch s := expr.(type) {
	case string:
		w.Assertions.Equal(expect, s)
	case data.Value:
		w.Assertions.Equal(expect, s.String())
	default:
		panic(fmt.Errorf(InvalidTestExpression, expr))
	}
}

// Number tests a Value for numeric equality
func (w *Wrapper) Number(expect float64, expr Any) {
	switch n := expr.(type) {
	case float64:
		w.Assertions.Equal(expect, float64(n))
	case int:
		w.Assertions.Equal(int64(expect), int64(n))
	case data.Number:
		w.Assertions.Equal(data.EqualTo, data.Float(expect).Cmp(n))
	default:
		panic(fmt.Errorf(InvalidTestExpression, expr))
	}
}

// Equal tests a Value for some kind of equality. Performs checks to do so
func (w *Wrapper) Equal(expect Any, expr Any) {
	switch typed := expect.(type) {
	case data.String:
		w.String(string(typed), expr)
	case data.Number:
		num := expr.(data.Number)
		w.Assertions.Equal(data.EqualTo, typed.Cmp(num))
	case data.Value:
		w.String(typed.String(), expr)
	default:
		w.Assertions.Equal(expect, expr)
	}
}

// True tests a Value for boolean true
func (w *Wrapper) True(expr Any) {
	if b, ok := expr.(data.Bool); ok {
		w.Assertions.True(bool(b))
		return
	}
	w.Assertions.True(expr.(bool))
}

// Truthy tests a Value for system-specific Truthy
func (w *Wrapper) Truthy(expr data.Value) {
	w.Assertions.True(data.Truthy(expr))
}

// False tests a Value for boolean false
func (w *Wrapper) False(expr Any) {
	if b, ok := expr.(data.Bool); ok {
		w.Assertions.False(bool(b))
		return
	}
	w.Assertions.False(expr.(bool))
}

// Falsey tests a Value for system-specific Falsey
func (w *Wrapper) Falsey(expr data.Value) {
	w.Assertions.False(data.Truthy(expr))
}

// Contains check if the expected string is in the provided Value
func (w *Wrapper) Contains(expect string, expr data.Value) {
	val := expr.String()
	w.Assertions.True(strings.Contains(val, expect))
}

// NotContains checks if the expected string is not in the provided Value
func (w *Wrapper) NotContains(expect string, expr data.Value) {
	val := expr.String()
	w.Assertions.False(strings.Contains(val, expect))
}

// Identical tests that two values are referentially identical
func (w *Wrapper) Identical(expect Any, expr Any) {
	w.Assertions.Equal(expect, expr)
}

// NotIdentical tests that two values are not referentially identical
func (w *Wrapper) NotIdentical(expect Any, expr Any) {
	w.Assertions.NotEqual(expect, expr)
}

// Compare tests if the Comparison of two Numbers is correct
func (w *Wrapper) Compare(c data.Comparison, l data.Number, r data.Number) {
	w.Assertions.Equal(c, l.Cmp(r))
}

// ExpectPanic is used with a defer to make sure an error was triggered
func (w *Wrapper) ExpectPanic(errStr string) {
	if rec := recover(); rec != nil {
		if re, ok := rec.(error); ok {
			recStr := re.Error()
			w.String(errStr, recStr)
			return
		}
	}
	panic("proper error not raised")
}

// ExpectNoPanic is sued with a defer to make sure no error was triggered
func (w *Wrapper) ExpectNoPanic() {
	rec := recover()
	w.Nil(rec)
}
