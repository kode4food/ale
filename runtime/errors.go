package runtime

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
)

type aleRuntimeError struct {
	message string
	wrapped error
}

// Error messages
const (
	ErrUnexpectedType = "got %s, expected %s"
)

var (
	firstCamel = regexp.MustCompile(`(.)([A-Z][a-z]+)`)
	restCamel  = regexp.MustCompile(`([a-z0-9])([A-Z])`)

	interfaceConversion = regexp.MustCompile(
		`^interface conversion: ` +
			`[^.]+[.](?P<got>[a-zA-Z0-9]+) ` +
			`is not [^.]+[.](?P<expected>[a-zA-Z0-9]+):.*$`,
	)
)

func AleRuntimeError(wrapped error, format string, a ...any) error {
	return &aleRuntimeError{
		message: fmt.Sprintf(format, a...),
		wrapped: wrapped,
	}
}

func (a *aleRuntimeError) Error() string {
	return a.message
}

func (a *aleRuntimeError) Unwrap() error {
	return a.wrapped
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

func normalizeTypeAssertionError(e *runtime.TypeAssertionError) error {
	if m := interfaceConversion.FindStringSubmatch(e.Error()); m != nil {
		return AleRuntimeError(e,
			ErrUnexpectedType, camelToWords(m[1]), camelToWords(m[2]),
		)
	}
	// Programmer error
	panic("could not normalize type assertion error")
}

func camelToWords(s string) string {
	res := firstCamel.ReplaceAllString(s, "${1} ${2}")
	res = restCamel.ReplaceAllString(res, "${1} ${2}")
	return strings.ToLower(res)
}
