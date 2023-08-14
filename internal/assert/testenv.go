package assert

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/core/bootstrap"
	"github.com/kode4food/ale/env"
)

var testEnv = bootstrap.DevNullEnvironment()

// GetTestEnvironment returns an immutable root testing Environment
func GetTestEnvironment() *env.Environment {
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
