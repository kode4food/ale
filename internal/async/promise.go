package async

import (
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/do"
	"github.com/kode4food/ale/internal/types"
)

type (
	// A Promise represents a Value that will eventually be resolved
	Promise struct {
		once     do.Action
		resolver data.Procedure
		result   any
		status   promiseStatus
	}

	promiseStatus int
)

const (
	promisePending promiseStatus = iota
	promiseResolved
	promiseFailed
)

var (
	PromiseType = types.MakeBasic("promise")

	promiseArityChecker = data.MakeFixedChecker(0)
)

// NewPromise instantiates a new Promise
func NewPromise(resolver data.Procedure) *Promise {
	return &Promise{
		once:     do.Once(),
		resolver: resolver,
		status:   promisePending,
	}
}

func (p *Promise) Call(...data.Value) data.Value {
	p.once(func() {
		defer func() {
			if rec := recover(); rec != nil {
				p.result = rec
				p.status = promiseFailed
			}
		}()
		p.result = p.resolver.Call()
		p.status = promiseResolved
	})

	if p.status == promiseFailed {
		panic(p.result)
	}
	return p.result.(data.Value)
}

func (p *Promise) CheckArity(c int) error {
	return promiseArityChecker(c)
}

func (p *Promise) IsResolved() bool {
	return p.status != promisePending
}

func (p *Promise) Type() types.Type {
	return PromiseType
}

func (p *Promise) Equal(v data.Value) bool {
	return p == v
}

func (p *Promise) String() string {
	return data.DumpString(p)
}
