package runtime

import (
	"fmt"
	"regexp"
	"runtime"
	"slices"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/internal/strings"
)

type aleRuntimeError struct {
	*data.Object
	message string
	wrapped error
}

// ErrUnexpectedType maps a Go interface conversion error to something that
// will make more sense to an Ale program
const ErrUnexpectedType = "got %s, expected %s"

var interfaceConversion = []*regexp.Regexp{
	regexp.MustCompile(
		`: [^.]+[.](?P<got>[a-zA-Z0-9]+) ` +
			`is not [^.]+[.](?P<expected>[a-zA-Z0-9]+).*$`,
	),
	regexp.MustCompile(
		` is [^.]+[.](?P<got>[a-zA-Z0-9]+), ` +
			`not [^.]+[.](?P<expected>[a-zA-Z0-9]+).*$`,
	),
	regexp.MustCompile(
		` is (?P<got>nil), ` +
			`not [^.]+[.](?P<expected>[a-zA-Z0-9]+).*$`),
}

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
	for _, re := range interfaceConversion {
		if m := re.FindStringSubmatch(e.Error()); m != nil {
			names := re.SubexpNames()
			expected := slices.Index(names, "expected")
			got := slices.Index(names, "got")
			return AleRuntimeError(e,
				ErrUnexpectedType,
				strings.CamelToWords(m[got]),
				strings.CamelToWords(m[expected]),
			)
		}
	}
	panic(debug.ProgrammerError("could not normalize type assertion error"))
}
