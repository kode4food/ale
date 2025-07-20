package builtin

import (
	"os"
	"regexp"
	"time"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/basics"
)

var envPairRegex = regexp.MustCompile("^(?P<Key>[^=]+)=(?P<Value>.*)$")

// Env returns an object containing the operating system's environment
func Env() ale.Value {
	var p data.Pairs
	for _, v := range os.Environ() {
		if e := envPairRegex.FindStringSubmatch(v); len(e) == 3 {
			p = append(p, data.NewCons(data.Keyword(e[1]), data.String(e[2])))
		}
	}
	return data.NewObject(p...)
}

// Args returns a vector containing the args passed to this program
func Args() ale.Value {
	return data.Vector(basics.Map(os.Args, func(v string) ale.Value {
		return data.String(v)
	}))
}

// CurrentTime returns the current system time in nanoseconds
var CurrentTime = data.MakeProcedure(func(...ale.Value) ale.Value {
	return data.Integer(time.Now().UnixNano())
}, 0)
