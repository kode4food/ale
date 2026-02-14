package data

import (
	"errors"
	"fmt"
)

const (
	// ErrFixedArity is raised when there are too few or many arguments
	// provided to a fixed ArityChecker
	ErrFixedArity = "expected %d arguments, got %d"

	// ErrMinimumArity is raised when there are too few arguments provided to a
	// minimum ArityChecker
	ErrMinimumArity = "expected at least %d arguments, got %d"

	// ErrRangedArity is raised when there are too few or many arguments
	// provided to a ranged ArityChecker
	ErrRangedArity = "expected between %d and %d arguments, got %d"

	// ErrTooManyArguments is raised when there are too many arguments provided
	// to a maximum ArityChecker
	ErrTooManyArguments = "too many arity check arguments"
)

// OrMore is the constant used when you want to tell MakeArityChecker to
// generate a minimum arity checker
const OrMore = -1

// MakeArityChecker produces an arity checker based on its parameters
func MakeArityChecker(arity ...int) (ArityChecker, error) {
	al := len(arity)
	switch {
	case al == 0:
		return CheckAnyArity, nil
	case al > 2:
		return nil, errors.New(ErrTooManyArguments)
	case al == 1 || arity[0] == arity[1]:
		return makeFixedChecker(arity[0]), nil
	case al == 2 && arity[1] == OrMore:
		return makeMinimumChecker(arity[0]), nil
	default:
		return makeRangedChecker(arity[0], arity[1]), nil
	}
}

// CheckAnyArity allows for any number of arguments
func CheckAnyArity(int) error {
	return nil
}

func makeFixedChecker(fixed int) ArityChecker {
	return func(count int) error {
		return CheckFixedArity(fixed, count)
	}
}

// CheckFixedArity allows for a fixed number of arguments
func CheckFixedArity(fixed, count int) error {
	if count != fixed {
		return fmt.Errorf(ErrFixedArity, fixed, count)
	}
	return nil
}

// MustCheckFixedArity is CheckFixedArity, but panics if the check fails
func MustCheckFixedArity(fixed, count int) {
	if err := CheckFixedArity(fixed, count); err != nil {
		panic(err)
	}
}

func makeMinimumChecker(min int) ArityChecker {
	return func(count int) error {
		return CheckMinimumArity(min, count)
	}
}

// CheckMinimumArity allows for a minimum number of arguments
func CheckMinimumArity(min, count int) error {
	if count < min {
		return fmt.Errorf(ErrMinimumArity, min, count)
	}
	return nil
}

func makeRangedChecker(min, max int) ArityChecker {
	return func(count int) error {
		return CheckRangedArity(min, max, count)
	}
}

// CheckRangedArity allows for a ranged number of arguments
func CheckRangedArity(min, max, count int) error {
	if count < min || count > max {
		return fmt.Errorf(ErrRangedArity, min, max, count)
	}
	return nil
}
