package bootstrap

import (
	"os"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/core/internal/builtin"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/do"
	"github.com/kode4food/ale/internal/stream"
	"github.com/kode4food/ale/macro"
)

type (
	bootstrap struct {
		environment *env.Environment
		macroMap    macroMap
		specialMap  specialMap
		funcMap     funcMap
	}

	macroMap   map[data.Local]macro.Call
	specialMap map[data.Local]encoder.Call
	funcMap    map[data.Local]data.Function
)

var (
	topLevelOnce = do.Once()
	topLevel     *env.Environment
	devNullOnce  = do.Once()
	devNull      *env.Environment
)

// Into sets up initial built-ins and assets. This call is useful if you're
// wiring up your own Environments. Otherwise, calls to TopLevelEnvironment and
// DevNullEnvironment will perform this action for you.
func Into(e *env.Environment) {
	b := &bootstrap{
		environment: e,
		macroMap:    macroMap{},
		specialMap:  specialMap{},
		funcMap:     funcMap{},
	}
	b.builtIns()
	b.assets()
}

// ProcessEnv binds *env* to the operating system's environment variables
func ProcessEnv(e *env.Environment) {
	ns := e.GetRoot()
	ns.Declare("*env*").Bind(builtin.Env())
}

// ProcessArgs binds *args* to the current Go app's command line arguments
func ProcessArgs(e *env.Environment) {
	ns := e.GetRoot()
	ns.Declare("*args*").Bind(builtin.Args())
}

// StandardIO binds *in*, *out*, and *err* to the operating system's standard
// input and output facilities
func StandardIO(e *env.Environment) {
	ns := e.GetRoot()
	ns.Declare("*in*").Bind(stream.NewReader(os.Stdin, stream.LineInput))
	ns.Declare("*out*").Bind(stream.NewWriter(os.Stdout, stream.StrOutput))
	ns.Declare("*err*").Bind(stream.NewWriter(os.Stderr, stream.StrOutput))
}

// DevNull binds *in*, *out*, and *err* to the operating system's bit bucket
// device (usually /dev/null)
func DevNull(e *env.Environment) {
	ns := e.GetRoot()
	devNull, _ := os.OpenFile(os.DevNull, os.O_RDWR, 0666)
	ns.Declare("*in*").Bind(stream.NewReader(devNull, stream.LineInput))
	ns.Declare("*out*").Bind(stream.NewWriter(devNull, stream.StrOutput))
	ns.Declare("*err*").Bind(stream.NewWriter(devNull, stream.StrOutput))
}

// TopLevelEnvironment configures an environment that could be used at the
// top-level of the system, such as the REPL. It has access to the *env*,
// *args*, and operating system's standard in/out/err file streams.
func TopLevelEnvironment() *env.Environment {
	topLevelOnce(func() {
		topLevel = env.NewEnvironment()
		ProcessEnv(topLevel)
		ProcessArgs(topLevel)
		StandardIO(topLevel)
		Into(topLevel)
	})
	res, err := topLevel.Snapshot()
	if err != nil {
		panic(err)
	}
	return res
}

// DevNullEnvironment configures a bootstrapped environment that is completely
// isolated from the top-level of the system. All I/O is rerouted to and from
// the operating system's bit bucket device (usually /dev/null)
func DevNullEnvironment() *env.Environment {
	devNullOnce(func() {
		devNull = env.NewEnvironment()
		DevNull(devNull)
		Into(devNull)
	})
	res, err := devNull.Snapshot()
	if err != nil {
		panic(err)
	}
	return res
}
