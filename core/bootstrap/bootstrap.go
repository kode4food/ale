package bootstrap

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/kode4food/ale/core/builtin"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/compiler"
	lang "github.com/kode4food/ale/internal/lang/env"
	"github.com/kode4food/ale/internal/stream"
	"github.com/kode4food/ale/internal/sync"
	"github.com/kode4food/ale/macro"
)

type (
	bootstrap struct {
		environment *env.Environment
		macroMap    macroMap
		specialMap  specialMap
		procMap     procMap
	}

	macroMap   map[data.Local]macro.Call
	specialMap map[data.Local]compiler.Call
	procMap    map[data.Local]data.Procedure
)

// Error messages
const (
	ErrCannotDeclareFS = "failed to declare filesystem: %w"
	ErrCannotBindFS    = "failed to bind filesystem: %w"
)

var (
	topLevelOnce = sync.Once()
	topLevel     *env.Environment
	devNullOnce  = sync.Once()
	devNull      *env.Environment
)

// Into sets up initial built-ins and populateAssets. Useful if you're wiring
// up your own Environments. Otherwise, calls to TopLevelEnvironment and
// DevNullEnvironment will perform this action for you.
func Into(e *env.Environment) {
	b := &bootstrap{
		environment: e,
		macroMap:    macroMap{},
		specialMap:  specialMap{},
		procMap:     procMap{},
	}
	b.populateDefiners()
	b.populateSpecialForms()
	b.populateBuiltins()
	b.populateAssets()
}

// ProcessEnv binds *env* to the operating system's environment variables
func ProcessEnv(e *env.Environment) {
	mustBindPublic(e.GetRoot(), lang.Env, builtin.Env())
}

// ProcessArgs binds *args* to the current Go app's command line arguments
func ProcessArgs(e *env.Environment) {
	mustBindPublic(e.GetRoot(), lang.Args, builtin.Args())
}

// StandardIO binds *in*, *out*, and *err* to the operating system's standard
// input and output facilities
func StandardIO(e *env.Environment) {
	ns := e.GetRoot()
	mustBindPublic(ns, lang.In, stream.NewReader(os.Stdin, stream.LineInput))
	mustBindPublic(ns, lang.Out, stream.NewWriter(os.Stdout, stream.StrOutput))
	mustBindPublic(ns, lang.Err, stream.NewWriter(os.Stderr, stream.StrOutput))
}

// DevNull binds *in*, *out*, and *err* to the operating system's bit bucket
// device (usually /dev/null)
func DevNull(e *env.Environment) {
	ns := e.GetRoot()
	devNull, _ := os.OpenFile(os.DevNull, os.O_RDWR, 0666)
	mustBindPublic(ns, lang.In, stream.NewReader(devNull, stream.LineInput))
	mustBindPublic(ns, lang.Out, stream.NewWriter(devNull, stream.StrOutput))
	mustBindPublic(ns, lang.Err, stream.NewWriter(devNull, stream.StrOutput))
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
	return topLevel.Snapshot()
}

// DevNullEnvironment configures a bootstrapped environment completely isolated
// from the top-level of the system. All I/O is rerouted to and from the
// operating system's bit bucket device (usually /dev/null)
func DevNullEnvironment() *env.Environment {
	devNullOnce(func() {
		devNull = env.NewEnvironment()
		DevNull(devNull)
		Into(devNull)
	})
	return devNull.Snapshot()
}

// BindFileSystem binds a file system to a Namespace to enable source includes
// and other file operations. The file system is private to the Namespace
func BindFileSystem(ns env.Namespace, f fs.FS) error {
	e, err := ns.Private(lang.FS)
	if err != nil {
		return fmt.Errorf(ErrCannotDeclareFS, err)
	}
	if err = e.Bind(stream.WrapFileSystem(f)); err != nil {
		return fmt.Errorf(ErrCannotBindFS, err)
	}
	return nil
}

// MustBindFileSystem binds a file system to a Namespace or panics if it can't
func MustBindFileSystem(ns env.Namespace, f fs.FS) {
	if err := BindFileSystem(ns, f); err != nil {
		panic(err)
	}
}
