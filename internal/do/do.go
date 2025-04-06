package do

import (
	"sync"
	"sync/atomic"
)

// Action is a callback interface for eventually triggering some action
type Action func(func())

const (
	doneNever uint32 = iota
	doneOnce
)

// Once creates a Do instance for performing an action only once
func Once() Action {
	var state atomic.Uint32
	var mutex sync.Mutex

	execLocked := func(f func()) {
		mutex.Lock()
		defer mutex.Unlock()

		if state.Load() == doneNever {
			defer state.Store(doneOnce)
			f()
		}
	}

	return func(f func()) {
		if state.Load() == doneNever {
			execLocked(f)
		}
	}
}

// Always returns a Do instance for always performing an action
func Always() Action {
	return func(f func()) {
		f()
	}
}

// Never returns a Do instance for never performing an action
func Never() Action {
	return func(func()) {
		// no-op
	}
}
