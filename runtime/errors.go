package runtime

import (
	"fmt"
	"regexp"
	"runtime"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/strings"
)

type aleRuntimeError struct {
	data.Object
	message string
	wrapped error
}

// Error messages
const (
	ErrUnexpectedType = "got %s, expected %s"
)

var interfaceConversion = regexp.MustCompile(
	`^interface conversion: ` +
		`[^.]+[.](?P<got>[a-zA-Z0-9]+) ` +
		`is not [^.]+[.](?P<expected>[a-zA-Z0-9]+):.*$`,
)

func AleRuntimeError(wrapped error, format string, a ...any) error {
	message := fmt.Sprintf(format, a...)
	object, _ := data.ValuesToObject(
		data.Keyword("message"), data.String(message),
		data.Keyword("wrapped"), data.String(wrapped.Error()),
	)
	return &aleRuntimeError{
		Object:  object,
		message: message,
		wrapped: wrapped,
	}
}

func NormalizeGoRuntimeErrors() {
	if rec := recover(); rec != nil {
		panic(NormalizeGoRuntimeError(rec))
	}
}

func NormalizeGoRuntimeError(value any) any {
	switch value := value.(type) {
	case *runtime.TypeAssertionError:
		return normalizeTypeAssertionError(value)
	default:
		return value
	}
}

func (a *aleRuntimeError) Error() string {
	return a.message
}

func (a *aleRuntimeError) String() string {
	return a.message
}

func (a *aleRuntimeError) Unwrap() error {
	return a.wrapped
}

func normalizeTypeAssertionError(e *runtime.TypeAssertionError) error {
	if m := interfaceConversion.FindStringSubmatch(e.Error()); m != nil {
		return AleRuntimeError(e,
			ErrUnexpectedType,
			strings.CamelToWords(m[1]),
			strings.CamelToWords(m[2]),
		)
	}
	// Programmer error
	panic("could not normalize type assertion error")
}
