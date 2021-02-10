package builtin

import (
	"os"
	"regexp"
	"time"

	"github.com/kode4food/ale/data"
)

var envPairRegex = regexp.MustCompile("^(?P<Key>[^=]+)=(?P<Value>.*)$")

// Env returns an object containing the operating system's environment
func Env() data.Value {
	var p data.Pairs
	for _, v := range os.Environ() {
		if e := envPairRegex.FindStringSubmatch(v); len(e) == 3 {
			p = append(p, data.NewCons(data.Keyword(e[1]), data.String(e[2])))
		}
	}
	return data.NewObject(p...)
}

// Args returns a vector containing the args passed to this program
func Args() data.Value {
	var r data.Values
	for _, v := range os.Args {
		r = append(r, data.String(v))
	}
	return data.NewVector(r...)
}

// CurrentTime returns the current system time in nanoseconds
var CurrentTime = data.Applicative(func(_ ...data.Value) data.Value {
	return data.Integer(time.Now().UnixNano())
}, 0)
