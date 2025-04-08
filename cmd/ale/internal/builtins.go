package internal

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/internal/compiler/special"
	"github.com/kode4food/ale/internal/console"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

func (r *REPL) registerBuiltIns() {
	r.registerBuiltIn("cls", data.MakeProcedure(cls, 0))
	r.registerBuiltIn("doc", doc)
	r.registerBuiltIn("debug", data.MakeProcedure(debugInfo, 0))
	r.registerBuiltIn("help", data.MakeProcedure(help, 0))
	r.registerBuiltIn("quit", data.MakeProcedure(shutdown, 0))
	r.registerBuiltIn("use", r.makeUse())
}

func (r *REPL) registerBuiltIn(n data.Local, v data.Value) {
	ns := r.getBuiltInsNamespace()
	_ = env.BindPublic(ns, n, v)
}

func (r *REPL) getBuiltInsNamespace() env.Namespace {
	return r.ns.Environment().GetRoot()
}

func (r *REPL) makeUse() data.Value {
	return special.Call(func(e encoder.Encoder, args ...data.Value) error {
		if err := data.CheckFixedArity(1, len(args)); err != nil {
			return err
		}
		n := args[0].(data.Local)
		old := r.ns
		r.ns = r.ns.Environment().GetQualified(n)
		if old != r.ns {
			fmt.Println()
		}
		return generate.Literal(e, nothing)
	})
}

func shutdown(...data.Value) data.Value {
	t := time.Now().UTC().UnixNano()
	rs := rand.NewSource(t)
	rg := rand.New(rs)
	idx := rg.Intn(len(farewells))
	fmt.Println(farewells[idx])
	os.Exit(0)
	return nothing
}

func debugInfo(...data.Value) data.Value {
	runtime.GC()
	fmt.Println("Number of goroutines: ", runtime.NumGoroutine())
	return nothing
}

func cls(...data.Value) data.Value {
	fmt.Println(console.Clear)
	return nothing
}
