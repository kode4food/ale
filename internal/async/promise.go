package async

import (
	"github.com/kode4food/ale/compiler/arity"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/do"
)

type (
	// Promise represents a Value that will eventually be resolved
	Promise interface {
		data.Caller
		IsResolved() bool
	}

	promiseStatus int

	promise struct {
		once     do.Do
		resolver data.Call
		result   interface{}
		status   promiseStatus
	}
)

const (
	promisePending promiseStatus = iota
	promiseResolved
	promiseFailed
)

var promiseArityChecker = arity.MakeFixedChecker(0)

// NewPromise instantiates a new Promise
func NewPromise(resolver data.Call) Promise {
	return &promise{
		once:     do.Once(),
		resolver: resolver,
		status:   promisePending,
	}
}

func (p *promise) Call() data.Call {
	return func(args ...data.Value) data.Value {
		p.once(func() {
			defer func() {
				if rec := recover(); rec != nil {
					p.result = rec
					p.status = promiseFailed
				}
			}()
			p.result = p.resolver()
			p.status = promiseResolved
		})

		if p.status == promiseFailed {
			panic(p.result)
		}
		return p.result.(data.Value)
	}
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

func (p *promise) Type() data.Name {
	return "promise"
}

func (p *promise) String() string {
	return data.DumpString(p)
}
