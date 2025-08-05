package sync

import (
	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/types"
)

type (
	// A Promise represents a Value that will eventually be resolved
	Promise struct {
		once     Action
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

var PromiseType = types.MakeBasic("promise")

// NewPromise instantiates a new Promise
func NewPromise(resolver data.Procedure) *Promise {
	return &Promise{
		once:     Once(),
		resolver: resolver,
		status:   promisePending,
	}
}

func (p *Promise) Call(...ale.Value) ale.Value {
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
	return p.result.(ale.Value)
}

func (p *Promise) CheckArity(c int) error {
	return data.CheckFixedArity(0, c)
}

func (p *Promise) IsResolved() bool {
	return p.status != promisePending
}

func (p *Promise) Type() ale.Type {
	return types.MakeLiteral(PromiseType, p)
}

func (p *Promise) Equal(other ale.Value) bool {
	return p == other
}

func (p *Promise) Get(key ale.Value) (ale.Value, bool) {
	return data.DumpMapped(p).Get(key)
}
