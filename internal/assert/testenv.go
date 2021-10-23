package assert

import (
	"sync"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/core/bootstrap"
	"github.com/kode4food/ale/env"
)

var (
	testEnv    = env.NewEnvironment()
	readyMutex sync.Mutex
	ready      bool
)

// GetTestEnvironment returns an immutable root testing Environment
func GetTestEnvironment() *env.Environment {
	readyMutex.Lock()
	defer readyMutex.Unlock()

	if !ready {
		bootstrap.Into(testEnv)
		ready = true
	}
	return testEnv
}

// GetTestNamespace returns a new anonymous namespace for testing purposes
func GetTestNamespace() env.Namespace {
	return GetTestEnvironment().GetAnonymous()
}

// GetTestEncoder returns a new Encoder, rooted at a test Namespace
func GetTestEncoder() encoder.Encoder {
	ns := GetTestNamespace()
	return encoder.NewEncoder(ns)
}
