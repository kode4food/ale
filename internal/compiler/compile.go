package compiler

import "gitlab.com/kode4food/ale/api"

// IsEvaluable returns whether or not the provided value is subject
// to further evaluation
func IsEvaluable(v api.Value) bool {
	switch v.(type) {
	case api.String:
		return false
	case api.Sequence, api.Symbol:
		return true
	default:
		return false
	}
}
