package internal

import (
	"fmt"
	"math/rand/v2"
	"os"
	"runtime"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/cmd/ale/internal/console"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/compiler"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
)

func (r *REPL) registerBuiltIns() {
	r.registerBuiltIn("cls", data.MakeProcedure(cls, 0))
	r.registerBuiltIn("doc", doc)
	r.registerBuiltIn("debug", data.MakeProcedure(debugInfo, 0))
	r.registerBuiltIn("help", data.MakeProcedure(help, 0))
	r.registerBuiltIn("quit", data.MakeProcedure(shutdown, 0))
	r.registerBuiltIn("use", r.makeUse())
}

func (r *REPL) registerBuiltIn(n data.Local, v ale.Value) {
	ns := r.getBuiltInsNamespace()
	_ = env.BindPublic(ns, n, v)
}

func (r *REPL) getBuiltInsNamespace() env.Namespace {
	return r.ns.Environment().GetRoot()
}

func (r *REPL) makeUse() ale.Value {
	return compiler.Call(func(e encoder.Encoder, args ...ale.Value) error {
		if err := data.CheckFixedArity(1, len(args)); err != nil {
			return err
		}
		n := args[0].(data.Local)
		old := r.ns
		r.ns = env.MustGetQualified(r.ns.Environment(), n)
		if old != r.ns {
			fmt.Println()
		}
		return generate.Literal(e, nothing)
	})
}

func shutdown(...ale.Value) ale.Value {
	idx := rand.IntN(len(farewells))
	fmt.Println(farewells[idx])
	os.Exit(0)
	return nothing
}

func debugInfo(...ale.Value) ale.Value {
	runtime.GC()
	fmt.Println("Number of goroutines: ", runtime.NumGoroutine())
	return nothing
}

func cls(...ale.Value) ale.Value {
	fmt.Println(console.Clear)
	return nothing
}
