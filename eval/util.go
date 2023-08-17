package eval

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
)

// Error messages
const (
	ErrUnexpectedType = "got %s, expected %s"
)

var (
	firstCamel = regexp.MustCompile(`(.)([A-Z][a-z]+)`)
	restCamel  = regexp.MustCompile(`([a-z0-9])([A-Z])`)

	intfConv = regexp.MustCompile(
		`^interface conversion: ` +
			`[^.]+[.](?P<got>[a-zA-Z0-9]+) ` +
			`is not [^.]+[.](?P<expected>[a-zA-Z0-9]+):.*$`,
	)
)

func NormalizeGoRuntimeErrors() {
	if rec := recover(); rec != nil {
		panic(NormalizeGoRuntimeError(rec))
	}
}

func NormalizeGoRuntimeError(rec interface{}) interface{} {
	switch rec := rec.(type) {
	case *runtime.TypeAssertionError:
		return normalizeTypeAssertionError(rec)
	default:
		return rec
	}
}

func normalizeTypeAssertionError(e *runtime.TypeAssertionError) error {
	if m := intfConv.FindStringSubmatch(e.Error()); m != nil {
		return fmt.Errorf(
			ErrUnexpectedType, camelToWords(m[1]), camelToWords(m[2]),
		)
	}
	panic("couldn't normalize type assertion error")
}

func camelToWords(s string) string {
	res := firstCamel.ReplaceAllString(s, "${1} ${2}")
	res = restCamel.ReplaceAllString(res, "${1} ${2}")
	return strings.ToLower(res)
}
