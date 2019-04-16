package bootstrap

import (
	"os"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/builtin"
	"gitlab.com/kode4food/ale/internal/compiler/encoder"
	"gitlab.com/kode4food/ale/internal/macro"
	"gitlab.com/kode4food/ale/internal/namespace"
	"gitlab.com/kode4food/ale/stdlib"
)

type (
	bootstrap struct {
		manager    *namespace.Manager
		macroMap   macroMap
		specialMap specialMap
		funcMap    funcMap
	}

	macroMap   map[api.Name]macro.Call
	specialMap map[api.Name]encoder.Call
	funcMap    map[api.Name]*api.Function
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
	ns.Bind("*env*", builtin.Env())
	ns.Bind("*args*", builtin.Args())
	ns.Bind("*in*", builtin.MakeReader(os.Stdin, stdlib.LineInput))
	ns.Bind("*out*", builtin.MakeWriter(os.Stdout, stdlib.StrOutput))
	ns.Bind("*err*", builtin.MakeWriter(os.Stderr, stdlib.StrOutput))
	return manager
}

// DevNullManager configures a manager that is completely isolated from
// the top-level of the system. All I/O is rerouted to and from /dev/null
func DevNullManager() *namespace.Manager {
	manager := namespace.NewManager()
	ns := manager.GetRoot()
	devNull, _ := os.Open(os.DevNull)
	ns.Bind("*in*", builtin.MakeReader(devNull, stdlib.LineInput))
	ns.Bind("*out*", builtin.MakeWriter(devNull, stdlib.StrOutput))
	ns.Bind("*err*", builtin.MakeWriter(devNull, stdlib.StrOutput))
	return manager
}
