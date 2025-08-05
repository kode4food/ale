package data_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/types"
)

func mustParseInteger(s string) data.Number {
	res, err := data.ParseInteger(s)
	if err != nil {
		panic(err)
	}
	return res
}

func mustParseFloat(s string) data.Number {
	res, err := data.ParseFloat(s)
	if err != nil {
		panic(err)
	}
	return res
}

func mustParseRatio(s string) data.Number {
	res, err := data.ParseRatio(s)
	if err != nil {
		panic(err)
	}
	return res
}

func TestParseFloat(t *testing.T) {
	as := assert.New(t)

	n1 := mustParseFloat("12.8")
	n2 := data.Float(12.8)
	as.Equal(n1, n2)

	_, err := data.ParseFloat("'splosion!")
	as.EqualError(err, fmt.Sprintf(data.ErrExpectedFloat, S(`'splosion!`)))
}

func TestParseInteger(t *testing.T) {
	as := assert.New(t)

	n1 := mustParseInteger("37")
	n2 := data.Integer(37)
	as.Equal(n1, n2)

	_, err := data.ParseInteger("'splosion!")
	as.EqualError(err, fmt.Sprintf(data.ErrExpectedInteger, S(`'splosion!`)))
}

func TestParseRatio(t *testing.T) {
	as := assert.New(t)

	n1 := mustParseRatio("1/2")
	n2 := mustParseFloat("0.5")
	as.True(n1.Cmp(n2) == data.EqualTo)

	_, err := data.ParseRatio("'splosion!")
	as.EqualError(err, fmt.Sprintf(data.ErrExpectedRatio, S(`'splosion!`)))
}

func TestEqualTo(t *testing.T) {
	as := assert.New(t)
	n1 := data.Integer(20)
	n2 := data.Float(20.0)
	n3 := data.Float(25.75)
	n4 := mustParseRatio("40/2")
	n5 := mustParseRatio("5/2")

	as.Compare(data.EqualTo, n1, n1)
	as.Compare(data.EqualTo, n3, n3)
	as.Compare(data.EqualTo, n5, n5)

	as.Compare(data.EqualTo, n1, n2)
	as.Compare(data.EqualTo, n2, n1)
	as.Compare(data.EqualTo, n2, n4)
	as.Compare(data.EqualTo, n1, n4)

	as.Compare(data.EqualTo, n5, mustParseRatio("5/2"))
}

func TestNumberEqual(t *testing.T) {
	as := assert.New(t)
	n1 := data.Integer(20)
	n2 := data.Float(20.0)
	n3 := data.Float(25.75)
	n4 := mustParseRatio("40/2")
	n5 := mustParseRatio("5/2")
	n6 := mustParseInteger("100000000000000000000000000000000000000000000")

	as.True(n1.Equal(n1))
	as.True(n1.Equal(n4))
	as.True(n3.Equal(n3))
	as.True(n5.Equal(n5))
	as.True(n5.Equal(mustParseRatio("5/2")))
	as.True(n6.Equal(n6))

	as.False(n1.Equal(n2))
	as.False(n2.Equal(n1))
	as.False(n2.Equal(n4))
	as.False(n6.Equal(n1))
}

func TestLessThan(t *testing.T) {
	as := assert.New(t)
	n1 := mustParseFloat("12.8")
	n2 := data.Float(12.9)
	n3 := data.Integer(20)
	n4 := data.Float(20.0)
	n5 := data.Float(25.75)
	n6 := data.Integer(25)
	n7 := mustParseRatio("3/4")

	as.Compare(data.LessThan, n1, n2)
	as.Compare(data.LessThan, n1, n3)
	as.Compare(data.LessThan, n2, n3)
	as.Compare(data.LessThan, n2, n4)
	as.Compare(data.LessThan, n3, n5)
	as.Compare(data.LessThan, n3, n6)
	as.Compare(data.LessThan, n7, n6)
}

func TestGreaterThan(t *testing.T) {
	as := assert.New(t)
	n1 := mustParseFloat("12.8")
	n2 := data.Float(12.9)
	n3 := data.Integer(20)
	n4 := data.Float(20.0)
	n5 := data.Float(25.75)
	n6 := data.Integer(25)
	n7 := mustParseInteger("100000000000000000000000000000000000000000000")

	as.Compare(data.GreaterThan, n2, n1)
	as.Compare(data.GreaterThan, n3, n1)
	as.Compare(data.GreaterThan, n3, n2)
	as.Compare(data.GreaterThan, n4, n2)
	as.Compare(data.GreaterThan, n5, n3)
	as.Compare(data.GreaterThan, n6, n3)
	as.Compare(data.GreaterThan, n7, n6)
}

func TestMultiplication(t *testing.T) {
	as := assert.New(t)
	n1 := data.Integer(20)
	n2 := data.Float(20.0)
	n3 := data.Integer(5)
	n4 := data.Float(5.0)
	n5 := data.Float(9.25)
	n6 := mustParseRatio("1/2")
	n7 := mustParseInteger("1000000000000000000000000000000000000000000")
	n8 := mustParseRatio("1/5")

	as.Number(100.0, n1.Mul(n4))
	as.Number(100.0, n2.Mul(n3))
	as.Number(100.0, n2.Mul(n4))
	as.Number(46.25, n4.Mul(n5))
	as.Number(100, n1.Mul(n3))
	as.String("2.5", n4.Mul(n6))
	as.String("5000000000000000000000000000000000000000000", n7.Mul(n3))
	as.String("1/10", n6.Mul(n8))
	as.String("2.5", n6.Mul(n4))
}

func TestDivision(t *testing.T) {
	as := assert.New(t)
	n1 := data.Integer(20)
	n2 := data.Float(20.0)
	n3 := data.Integer(5)
	n4 := data.Float(5.0)
	n5 := mustParseRatio("1/2")
	n6 := mustParseInteger("1000000000000000000000000000000000000000000")

	as.Number(4.0, n1.Div(n4))
	as.Number(4.0, n2.Div(n3))
	as.Number(4.0, n2.Div(n4))
	as.Number(4, n1.Div(n3))
	as.String("40", n1.Div(n5))
	as.String("1/10", n5.Div(n3))
	as.String("50000000000000000000000000000000000000000", n6.Div(n1))
}

func TestRemainder(t *testing.T) {
	as := assert.New(t)
	n1 := data.Integer(5)
	n2 := data.Float(5.0)
	n3 := data.Integer(7)
	n4 := data.Float(7.0)
	n5 := mustParseRatio("14/2")
	n6 := mustParseInteger("1000000000000000000000000000000000000000000")
	n7 := mustParseRatio("15/2")
	n8 := mustParseRatio("7/2")
	n9 := data.Integer(-20)

	as.Number(2.0, n3.Mod(n2))
	as.Number(2.0, n4.Mod(n1))
	as.Number(2.0, n4.Mod(n2))
	as.Number(2, n3.Mod(n1))
	as.String("2", n5.Mod(n1))
	as.String("1", n6.Mod(n3))
	as.String("15/7", n7.Div(n8))
	as.String("1/2", n7.Mod(n8))
	as.String("1.0", n9.Mod(n4))
	as.String("1", n9.Mod(n3))
}

func TestAddition(t *testing.T) {
	as := assert.New(t)
	n1 := data.Integer(20)
	n2 := data.Integer(5)
	n3 := data.Float(9.25)
	n4 := data.Float(7.0)
	n5 := mustParseRatio("14/2")
	n6 := mustParseInteger("1000000000000000000000000000000000000000000")

	as.Number(16.25, n3.Add(n4))
	as.Number(29.25, n1.Add(n3))
	as.Number(29.25, n3.Add(n1))
	as.Number(25, n1.Add(n2))
	as.Number(25, n2.Add(n1))
	as.String("14.0", n4.Add(n5))
	as.String("1000000000000000000000000000000000000000005", n6.Add(n2))
}

func TestSubtraction(t *testing.T) {
	as := assert.New(t)
	n1 := data.Integer(20)
	n2 := data.Float(20.0)
	n3 := data.Integer(5)
	n4 := data.Float(5.0)
	n5 := data.Float(9.25)
	n6 := data.Float(7.0)
	n7 := mustParseRatio("14/2")
	n8 := mustParseInteger("1000000000000000000000000000000000000000000")
	n9 := mustParseRatio("1/5")

	as.Number(-15.0, n4.Sub(n1))
	as.Number(15.0, n1.Sub(n4))
	as.Number(15.0, n2.Sub(n3))
	as.Number(2.25, n5.Sub(n6))
	as.Number(15, n1.Sub(n3))
	as.Number(-9.05, n9.Sub(n5))
	as.String("2.0", n7.Sub(n4))
	as.String("34/5", n7.Sub(n9))
	as.String("999999999999999999999999999999999999999980", n8.Sub(n1))
}

func TestInfiniteNumbers(t *testing.T) {
	as := assert.New(t)

	as.False(data.Integer(98).IsPosInf())
	as.False(data.Integer(0).IsNegInf())

	posInf := data.Float(1).Div(data.Float(0))
	negInf := data.Float(-1).Div(data.Float(0))

	as.True(posInf.IsPosInf())
	as.False(posInf.IsNegInf())
	as.True(negInf.IsNegInf())
	as.False(negInf.IsPosInf())

	as.Compare(data.GreaterThan, posInf, data.Integer(1))
	as.Compare(data.LessThan, negInf, data.Integer(1))
	as.Compare(data.LessThan, data.Integer(1), posInf)
	as.Compare(data.GreaterThan, data.Integer(1), negInf)

	as.Compare(data.LessThan, negInf, data.Float(1))
	as.Compare(data.GreaterThan, posInf, data.Float(1))
	as.Compare(data.LessThan, data.Float(1), posInf)
	as.Compare(data.GreaterThan, data.Float(1), negInf)
}

func TestNonNumbers(t *testing.T) {
	as := assert.New(t)

	nan := data.Float(math.Log(-1.0))
	as.True(nan.IsNaN())

	as.Compare(data.Incomparable, data.Float(1), nan)
	as.Compare(data.Incomparable, nan, data.Float(1))
	as.Compare(data.Incomparable, data.Integer(1), nan)
	as.Compare(data.Incomparable, nan, data.Integer(1))
}

func TestStringifyNumbers(t *testing.T) {
	as := assert.New(t)
	n1 := mustParseFloat("12.8")
	n2 := data.Float(12.9)
	n3 := data.Float(20)

	as.String("12.8", n1)
	as.String("12.9", n2)
	as.String("20.0", n3)
}

func TestPurify(t *testing.T) {
	as := assert.New(t)

	n1 := mustParseFloat("0.5")
	n2 := mustParseInteger("99999999999999999999999999999999999999999")
	n3 := mustParseInteger("1")
	n4 := mustParseRatio("1/2")

	as.String("1.5", n1.Add(n3))
	as.String("1.0", n1.Add(n4))
	as.String("1e+41", n1.Add(n2))
	as.String("1e+41", n2.Add(n1))
	as.String("100000000000000000000000000000000000000000", n2.Add(n3))
	as.String("199999999999999999999999999999999999999999/2", n2.Add(n4))
	as.String("1.5", n3.Add(n1))
	as.String("100000000000000000000000000000000000000000", n3.Add(n2))
	as.String("3/2", n3.Add(n4))
	as.String("1.0", n4.Add(n1))
	as.String("199999999999999999999999999999999999999999/2", n4.Add(n2))
	as.String("3/2", n4.Add(n3))
}

func TestIntegerPromotion(t *testing.T) {
	as := assert.New(t)

	i1 := data.Integer(math.MaxInt64)
	i2 := data.Integer(math.MinInt64)

	r1, ok := i1.Add(data.Integer(1)).(*data.BigInt)
	as.True(ok)
	as.String("9223372036854775808", r1)

	r2, ok := i1.Add(data.Integer(0)).(data.Integer)
	as.True(ok)
	as.String("9223372036854775807", r2)

	r3, ok := i2.Sub(data.Integer(1)).(*data.BigInt)
	as.True(ok)
	as.String("-9223372036854775809", r3)

	r4, ok := i2.Sub(data.Integer(0)).(data.Integer)
	as.True(ok)
	as.String("-9223372036854775808", r4)

	r5, ok := i1.Mul(data.Integer(2)).(*data.BigInt)
	as.True(ok)
	as.String("18446744073709551614", r5)

	r6, ok := i1.Mul(data.Integer(1)).(data.Integer)
	as.True(ok)
	as.String("9223372036854775807", r6)
}

func TestIntegerDemotion(t *testing.T) {
	as := assert.New(t)

	i1 := data.Integer(math.MaxInt64)
	i2 := data.Integer(math.MinInt64)

	r1, ok := i1.Add(data.Integer(1)).(*data.BigInt)
	as.True(ok)
	as.String("9223372036854775808", r1)

	r2, ok := r1.Sub(data.Integer(1)).(data.Integer)
	as.True(ok)
	as.String("9223372036854775807", r2)

	r3, ok := i2.Sub(data.Integer(1)).(*data.BigInt)
	as.True(ok)
	as.String("-9223372036854775809", r3)

	r4, ok := r3.Add(data.Integer(1)).(data.Integer)
	as.True(ok)
	as.String("-9223372036854775808", r4)
}

func TestNumberHashing(t *testing.T) {
	as := assert.New(t)

	n1 := mustParseFloat("0.5")
	n2 := mustParseInteger("99999999999999999999999999999999999999999")
	n3 := mustParseInteger("1")
	n4 := mustParseRatio("1/2")

	o := O(
		C(n1, data.True),
		C(n2, data.True),
		C(n3, data.True),
		C(n4, data.True),
	)

	v, ok := o.Get(n1)
	as.True(ok)
	as.True(v)

	v, ok = o.Get(n2)
	as.True(ok)
	as.True(v)

	v, ok = o.Get(n3)
	as.True(ok)
	as.True(v)

	v, ok = o.Get(n4)
	as.True(ok)
	as.True(v)
}

func testNumberChecks(as *assert.Wrapper, n data.Number) {
	as.False(n.IsNaN())
	as.False(n.IsPosInf())
	as.False(n.IsNegInf())
	as.True(types.BasicNumber.Accepts(n.(ale.Typed).Type()))
}

func TestNumberChecks(t *testing.T) {
	as := assert.New(t)

	n1 := mustParseFloat("0.5")
	n2 := mustParseInteger("99999999999999999999999999999999999999999")
	n3 := mustParseInteger("1")
	n4 := mustParseRatio("1/2")

	testNumberChecks(as, n1)
	testNumberChecks(as, n2)
	testNumberChecks(as, n3)
	testNumberChecks(as, n4)
}
