package builtin

import "gitlab.com/kode4food/ale/api"

// Add returns the sum of the provided numbers
func Add(args ...api.Value) api.Value {
	var res api.Number = api.Integer(0)
	for _, n := range args {
		res = res.Add(n.(api.Number))
	}
	return res
}

// Sub will subtract one number from the previous, in turn
func Sub(args ...api.Value) api.Value {
	var res = args[0].(api.Number)
	for _, n := range args[1:] {
		res = res.Sub(n.(api.Number))
	}
	return res
}

// Mul will generate the product of all the provided numbers
func Mul(args ...api.Value) api.Value {
	var res api.Number = api.Integer(1)
	for _, n := range args {
		res = res.Mul(n.(api.Number))
	}
	return res
}

// Div will divide one number by the next, in turn
func Div(args ...api.Value) api.Value {
	var res = args[0].(api.Number)
	for _, n := range args[1:] {
		res = res.Div(n.(api.Number))
	}
	return res
}

// Mod will produce the remainder of dividing one number by the next, in turn
func Mod(args ...api.Value) api.Value {
	var res = args[0].(api.Number)
	for _, n := range args[1:] {
		res = res.Mod(n.(api.Number))
	}
	return res
}

// Eq returns whether or not the provided numbers are equal
func Eq(args ...api.Value) api.Value {
	var res = args[0].(api.Number)
	for _, n := range args[1:] {
		if res.Cmp(n.(api.Number)) != api.EqualTo {
			return api.False
		}
	}
	return api.True
}

// Neq returns true any of the numbers is not equal to the others
func Neq(args ...api.Value) api.Value {
	if Eq(args...) == api.True {
		return api.False
	}
	return api.True
}

// Gt returns true if each number is greater than the previous
func Gt(args ...api.Value) api.Value {
	var l = args[0].(api.Number)
	for _, v := range args[1:] {
		r := v.(api.Number)
		if l.Cmp(r) != api.GreaterThan {
			return api.False
		}
		l = r
	}
	return api.True
}

// Gte returns true if each number is greater than or equal to the previous
func Gte(args ...api.Value) api.Value {
	var l = args[0].(api.Number)
	for _, v := range args[1:] {
		r := v.(api.Number)
		cmp := l.Cmp(r)
		if !(cmp == api.GreaterThan || cmp == api.EqualTo) {
			return api.False
		}
		l = r
	}
	return api.True
}

// Lt returns true if each number is less than the previous
func Lt(args ...api.Value) api.Value {
	var l = args[0].(api.Number)
	for _, v := range args[1:] {
		r := v.(api.Number)
		if l.Cmp(r) != api.LessThan {
			return api.False
		}
		l = r
	}
	return api.True
}

// Lte returns true if each number is less than or equal to the previous
func Lte(args ...api.Value) api.Value {
	var l = args[0].(api.Number)
	for _, v := range args[1:] {
		r := v.(api.Number)
		cmp := l.Cmp(r)
		if !(cmp == api.LessThan || cmp == api.EqualTo) {
			return api.False
		}
		l = r
	}
	return api.True
}

// IsPosInf returns true if the provided number represents positive infinity
func IsPosInf(args ...api.Value) api.Value {
	if num, ok := args[0].(api.Number); ok {
		return api.Bool(num.IsPosInf())
	}
	return api.False
}

// IsNegInf returns true if the provided number represents negative infinity
func IsNegInf(args ...api.Value) api.Value {
	if num, ok := args[0].(api.Number); ok {
		return api.Bool(num.IsNegInf())
	}
	return api.False
}

// IsNaN returns true if the provided value is not a number
func IsNaN(args ...api.Value) api.Value {
	if num, ok := args[0].(api.Number); ok {
		return api.Bool(num.IsNaN())
	}
	return api.False
}
