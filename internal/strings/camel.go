package strings

import (
	"regexp"
	"strings"
)

var (
	firstCamel = regexp.MustCompile(`(.)([A-Z][a-z]+)`)
	restCamel  = regexp.MustCompile(`([a-z0-9])([A-Z])`)
)

func CamelToWords(s string) string {
	res := firstCamel.ReplaceAllString(s, "${1} ${2}")
	res = restCamel.ReplaceAllString(res, "${1} ${2}")
	return strings.ToLower(res)
}

func CamelToSnake(s string) string {
	res := firstCamel.ReplaceAllString(s, "${1}-${2}")
	res = restCamel.ReplaceAllString(res, "${1}-${2}")
	return strings.ToLower(res)
}
