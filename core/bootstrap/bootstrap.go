package bootstrap

import (
	"os"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/core/internal/builtin"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/stream"
	"github.com/kode4food/ale/macro"
	"github.com/kode4food/ale/namespace"
)

type (
	bootstrap struct {
		manager    *namespace.Manager
		macroMap   macroMap
		specialMap specialMap
		funcMap    funcMap
	}

	macroMap   map[data.Name]macro.Call
	specialMap map[data.Name]encoder.Call
	funcMap    map[data.Name]data.Function
)

// Into sets up initial built-ins and assets
func Into(manager *namespace.Manager) {
	b := &bootstrap{
		manager:    manager,
		macroMap:   macroMap{},
		specialMap: specialMap{},
		funcMap:    funcMap{},
	}
	b.builtIns()
	b.assets()
}

// TopLevelManager configures a manager that could be used at the top-level
// of the system, such as the REPL. It has access to the *env*, *args*, and
// standard in/out/err file streams.
func TopLevelManager() *namespace.Manager {
	manager := namespace.NewManager()
	ns := manager.GetRoot()
	ns.Declare("*env*").Bind(builtin.Env())
	ns.Declare("*args*").Bind(builtin.Args())
	ns.Declare("*in*").Bind(builtin.MakeReader(os.Stdin, stream.LineInput))
	ns.Declare("*out*").Bind(builtin.MakeWriter(os.Stdout, stream.StrOutput))
	ns.Declare("*err*").Bind(builtin.MakeWriter(os.Stderr, stream.StrOutput))
	return manager
}

// DevNullManager configures a manager that is completely isolated from
// the top-level of the system. All I/O is rerouted to and from /dev/null
func DevNullManager() *namespace.Manager {
	manager := namespace.NewManager()
	ns := manager.GetRoot()
	devNull, _ := os.OpenFile(os.DevNull, os.O_RDWR, 0666)
	ns.Declare("*in*").Bind(builtin.MakeReader(devNull, stream.LineInput))
	ns.Declare("*out*").Bind(builtin.MakeWriter(devNull, stream.StrOutput))
	ns.Declare("*err*").Bind(builtin.MakeWriter(devNull, stream.StrOutput))
	return manager
}
