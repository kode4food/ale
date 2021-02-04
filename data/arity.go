package data

import (
	"errors"
	"fmt"
)

// Error messages
const (
	ErrFixedArity       = "expected %d arguments, got %d"
	ErrMinimumArity     = "expected at least %d arguments, got %d"
	ErrRangedArity      = "expected between %d and %d arguments, got %d"
	ErrTooManyArguments = "too many arity check arguments"
)

// OrMore is the constant used when you want to tell MakeChecker
// to generate a minimum arity checker
const OrMore = -1

// MakeChecker produces an arity checker based on its parameters
func MakeChecker(arity ...int) ArityChecker {
	al := len(arity)
	switch {
	case al == 0:
		return nil
	case al > 2:
		panic(errors.New(ErrTooManyArguments))
	case al == 1 || arity[0] == arity[1]:
		return MakeFixedChecker(arity[0])
	case al == 2 && arity[1] == OrMore:
		return MakeMinimumChecker(arity[0])
	default:
		return MakeRangedChecker(arity[0], arity[1])
	}
}

// AssertFixed explodes if a fixed arity check fails
func AssertFixed(fixed, count int) int {
	if err := MakeFixedChecker(fixed)(count); err != nil {
		panic(err)
	}
	return count
}

// MakeFixedChecker generates a fixed arity checker
func MakeFixedChecker(fixed int) ArityChecker {
	return func(count int) error {
		if count != fixed {
			return fmt.Errorf(ErrFixedArity, fixed, count)
		}
		return nil
	}
}

// AssertMinimum explodes if a fixed arity check fails
func AssertMinimum(min, count int) int {
	if err := MakeMinimumChecker(min)(count); err != nil {
		panic(err)
	}
	return count
}

// MakeMinimumChecker generates a minimum arity checker
func MakeMinimumChecker(min int) ArityChecker {
	return func(count int) error {
		if count < min {
			return fmt.Errorf(ErrMinimumArity, min, count)
		}
		return nil
	}
}

// AssertRanged explodes if a fixed arity check fails
func AssertRanged(min, max, count int) int {
	if err := MakeRangedChecker(min, max)(count); err != nil {
		panic(err)
	}
	return count
}

// MakeRangedChecker generates a ranged arity checker
func MakeRangedChecker(min, max int) ArityChecker {
	return func(count int) error {
		if count < min || count > max {
			return fmt.Errorf(ErrRangedArity, min, max, count)
		}
		return nil
	}
}

func checkRangedArity(min, max, count int) error {
	return MakeRangedChecker(min, max)(count)
}
