package assert

import (
	"sync"

	"github.com/kode4food/ale/core/bootstrap"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/compiler/encoder"
)

var (
	testOnce sync.Once
	testEnv  *env.Environment
)

// GetTestEnvironment returns an immutable root testing Environment
func GetTestEnvironment() *env.Environment {
	testOnce.Do(func() {
		testEnv = bootstrap.DevNullEnvironment()
	})
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
