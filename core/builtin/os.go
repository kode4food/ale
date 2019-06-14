package builtin

import (
	"os"
	"regexp"
	"time"

	"gitlab.com/kode4food/ale/data"
)

var envPairRegex = regexp.MustCompile("^(?P<Key>[^=]+)=(?P<Value>.*)$")

// Env returns an associative containing the operating system's environment
func Env() data.Value {
	r := data.Values{}
	for _, v := range os.Environ() {
		if e := envPairRegex.FindStringSubmatch(v); len(e) == 3 {
			r = append(r, data.Keyword(e[1]), data.String(e[2]))
		}
	}
	return data.NewObject(r...)
}

// Args returns a vector containing the args passed to this program
func Args() data.Value {
	r := data.EmptyVector
	for _, v := range os.Args {
		r = append(r, data.String(v))
	}
	return r
}

// CurrentTime returns the current system time in nanoseconds
func CurrentTime(_ ...data.Value) data.Value {
	return data.Integer(time.Now().UnixNano())
}
