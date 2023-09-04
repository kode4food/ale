package async

import (
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/do"
)

type (
	// Promise represents a Value that will eventually be resolved
	Promise interface {
		data.Function
		IsResolved() bool
	}

	promiseStatus int

	promise struct {
		once     do.Action
		resolver data.Function
		result   any
		status   promiseStatus
	}
)

const (
	promisePending promiseStatus = iota
	promiseResolved
	promiseFailed
)

var promiseArityChecker = data.MakeFixedChecker(0)

// NewPromise instantiates a new Promise
func NewPromise(resolver data.Function) Promise {
	return &promise{
		once:     do.Once(),
		resolver: resolver,
		status:   promisePending,
	}
}

func (p *promise) Call(...data.Value) data.Value {
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

func (p *promise) Convention() data.Convention {
	return data.ApplicativeCall
}

func (p *promise) CheckArity(c int) error {
	return promiseArityChecker(c)
}

func (p *promise) IsResolved() bool {
	return p.status != promisePending
}

func (p *promise) Type() data.Local {
	return "promise"
}

func (p *promise) Equal(v data.Value) bool {
	return p == v
}

func (p *promise) String() string {
	return data.DumpString(p)
}
