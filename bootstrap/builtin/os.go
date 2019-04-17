package builtin

import (
	"os"
	"regexp"
	"time"

	"gitlab.com/kode4food/ale/api"
)

var envPairRegex = regexp.MustCompile("^(?P<Key>[^=]+)=(?P<Value>.*)$")

// Env returns an associative containing the operating system's environment
func Env() api.Value {
	var r []api.Vector
	for _, v := range os.Environ() {
		if e := envPairRegex.FindStringSubmatch(v); len(e) == 3 {
			r = append(r, api.Vector{
				api.Keyword(api.Name(e[1])),
				api.String(e[2]),
			})
		}
	}
	return api.NewAssociative(r...)
}

// Args returns a vector containing the args passed to this program
func Args() api.Value {
	r := api.EmptyVector
	for _, v := range os.Args {
		r = append(r, api.String(v))
	}
	return r
}

// CurrentTime returns the current system time in nanoseconds
func CurrentTime(_ ...api.Value) api.Value {
	return api.Integer(time.Now().UnixNano())
}
