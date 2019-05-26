package builtin

import "gitlab.com/kode4food/ale/data"

// Add returns the sum of the provided numbers
func Add(args ...data.Value) data.Value {
	if len(args) > 0 {
		var res = args[0].(data.Number)
		for _, n := range args[1:] {
			res = res.Add(n.(data.Number))
		}
		return res
	}
	return data.Integer(0)
}

// Sub will subtract one number from the previous, in turn
func Sub(args ...data.Value) data.Value {
	if len(args) > 1 {
		var res = args[0].(data.Number)
		var rest = args[1:]
		for _, n := range rest {
			res = res.Sub(n.(data.Number))
		}
		return res
	}
	return data.Integer(-1).Mul(args[0].(data.Number))
}

// Mul will generate the product of all the provided numbers
func Mul(args ...data.Value) data.Value {
	var res data.Number = data.Integer(1)
	for _, n := range args {
		res = res.Mul(n.(data.Number))
	}
	return res
}

// Div will divide one number by the next, in turn
func Div(args ...data.Value) data.Value {
	var res = args[0].(data.Number)
	for _, n := range args[1:] {
		res = res.Div(n.(data.Number))
	}
	return res
}

// Mod will produce the remainder of dividing one number by the next, in turn
func Mod(args ...data.Value) data.Value {
	var res = args[0].(data.Number)
	for _, n := range args[1:] {
		res = res.Mod(n.(data.Number))
	}
	return res
}

// Eq returns whether or not the provided numbers are equal
func Eq(args ...data.Value) data.Value {
	var res = args[0].(data.Number)
	for _, n := range args[1:] {
		if res.Cmp(n.(data.Number)) != data.EqualTo {
			return data.False
		}
	}
	return data.True
}

// Neq returns true if any of the numbers is not equal to the others
func Neq(args ...data.Value) data.Value {
	if Eq(args...) == data.True {
		return data.False
	}
	return data.True
}

// Gt returns true if each number is greater than the previous
func Gt(args ...data.Value) data.Value {
	var l = args[0].(data.Number)
	for _, v := range args[1:] {
		r := v.(data.Number)
		if l.Cmp(r) != data.GreaterThan {
			return data.False
		}
		l = r
	}
	return data.True
}

// Gte returns true if each number is greater than or equal to the previous
func Gte(args ...data.Value) data.Value {
	var l = args[0].(data.Number)
	for _, v := range args[1:] {
		r := v.(data.Number)
		cmp := l.Cmp(r)
		if !(cmp == data.GreaterThan || cmp == data.EqualTo) {
			return data.False
		}
		l = r
	}
	return data.True
}

// Lt returns true if each number is less than the previous
func Lt(args ...data.Value) data.Value {
	var l = args[0].(data.Number)
	for _, v := range args[1:] {
		r := v.(data.Number)
		if l.Cmp(r) != data.LessThan {
			return data.False
		}
		l = r
	}
	return data.True
}

// Lte returns true if each number is less than or equal to the previous
func Lte(args ...data.Value) data.Value {
	var l = args[0].(data.Number)
	for _, v := range args[1:] {
		r := v.(data.Number)
		cmp := l.Cmp(r)
		if !(cmp == data.LessThan || cmp == data.EqualTo) {
			return data.False
		}
		l = r
	}
	return data.True
}

// IsPosInf returns true if the provided number represents positive infinity
func IsPosInf(args ...data.Value) data.Value {
	if num, ok := args[0].(data.Number); ok {
		return data.Bool(num.IsPosInf())
	}
	return data.False
}

// IsNegInf returns true if the provided number represents negative infinity
func IsNegInf(args ...data.Value) data.Value {
	if num, ok := args[0].(data.Number); ok {
		return data.Bool(num.IsNegInf())
	}
	return data.False
}

// IsNaN returns true if the provided value is not a number
func IsNaN(args ...data.Value) data.Value {
	if num, ok := args[0].(data.Number); ok {
		return data.Bool(num.IsNaN())
	}
	return data.False
}
