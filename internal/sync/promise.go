package sync

import (
	"github.com/kode4food/ale/internal/types"
	"github.com/kode4food/ale/pkg/data"
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
	return data.CheckFixedArity(0, c)
}

func (p *Promise) IsResolved() bool {
	return p.status != promisePending
}

func (p *Promise) Type() types.Type {
	return PromiseType
}

func (p *Promise) Equal(other data.Value) bool {
	return p == other
}

func (p *Promise) Get(key data.Value) (data.Value, bool) {
	return data.DumpMapped(p).Get(key)
}
