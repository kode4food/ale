package data_test

import (
	"fmt"
	"math"
	"testing"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestParseFloat(t *testing.T) {
	as := assert.New(t)

	n1 := data.ParseFloat("12.8")
	n2 := F(12.8)
	as.Equal(n1, n2)

	defer as.ExpectPanic(fmt.Sprintf(data.ExpectedFloat, S(`'splosion!`)))
	data.ParseFloat("'splosion!")
}

func TestParseInteger(t *testing.T) {
	as := assert.New(t)

	n3 := data.ParseInteger("37")
	n4 := I(37)
	as.Equal(n3, n4)

	defer as.ExpectPanic(fmt.Sprintf(data.ExpectedInteger, S(`'splosion!`)))
	data.ParseInteger("'splosion!")
}

func TestEqualTo(t *testing.T) {
	as := assert.New(t)
	n1 := I(20)
	n2 := F(20.0)
	n5 := F(25.75)

	as.Compare(data.EqualTo, n1, n2)
	as.Compare(data.EqualTo, n2, n1)
	as.Compare(data.EqualTo, n1, n1)
	as.Compare(data.EqualTo, n5, n5)
}

func TestLessThan(t *testing.T) {
	as := assert.New(t)
	n1 := data.ParseFloat("12.8")
	n2 := F(12.9)
	n3 := I(20)
	n4 := F(20.0)
	n5 := F(25.75)
	n6 := I(25)

	as.Compare(data.LessThan, n1, n2)
	as.Compare(data.LessThan, n1, n3)
	as.Compare(data.LessThan, n2, n3)
	as.Compare(data.LessThan, n2, n4)
	as.Compare(data.LessThan, n3, n5)
	as.Compare(data.LessThan, n3, n6)
}

func TestGreaterThan(t *testing.T) {
	as := assert.New(t)
	n1 := data.ParseFloat("12.8")
	n2 := F(12.9)
	n3 := I(20)
	n4 := F(20.0)
	n5 := F(25.75)
	n6 := I(25)

	as.Compare(data.GreaterThan, n2, n1)
	as.Compare(data.GreaterThan, n3, n1)
	as.Compare(data.GreaterThan, n3, n2)
	as.Compare(data.GreaterThan, n4, n2)
	as.Compare(data.GreaterThan, n5, n3)
	as.Compare(data.GreaterThan, n6, n3)
}

func TestMultiplication(t *testing.T) {
	as := assert.New(t)
	n1 := I(20)
	n2 := F(20.0)
	n3 := I(5)
	n4 := F(5.0)
	n5 := F(9.25)

	as.Float(100.0, n1.Mul(n4))
	as.Float(100.0, n2.Mul(n3))
	as.Float(100.0, n2.Mul(n4))
	as.Float(46.25, n4.Mul(n5))
	as.Integer(100, n1.Mul(n3))
}

func TestDivision(t *testing.T) {
	as := assert.New(t)
	n1 := I(20)
	n2 := F(20.0)
	n3 := I(5)
	n4 := F(5.0)

	as.Float(4.0, n1.Div(n4))
	as.Float(4.0, n2.Div(n3))
	as.Float(4.0, n2.Div(n4))
	as.Integer(4, n1.Div(n3))
}

func TestRemainder(t *testing.T) {
	as := assert.New(t)
	n1 := I(5)
	n2 := F(5.0)
	n3 := I(7)
	n4 := F(7.0)

	as.Float(2.0, n3.Mod(n2))
	as.Float(2.0, n4.Mod(n1))
	as.Float(2.0, n4.Mod(n2))
	as.Integer(2, n3.Mod(n1))
}

func TestAddition(t *testing.T) {
	as := assert.New(t)
	n1 := I(20)
	n2 := I(5)
	n3 := F(9.25)
	n4 := F(7.0)

	as.Float(16.25, n3.Add(n4))
	as.Float(29.25, n1.Add(n3))
	as.Float(29.25, n3.Add(n1))
	as.Integer(25, n1.Add(n2))
	as.Integer(25, n2.Add(n1))
}

func TestSubtraction(t *testing.T) {
	as := assert.New(t)
	n1 := I(20)
	n2 := F(20.0)
	n3 := I(5)
	n4 := F(5.0)
	n5 := F(9.25)
	n6 := F(7.0)

	as.Float(-15.0, n4.Sub(n1))
	as.Float(15.0, n1.Sub(n4))
	as.Float(15.0, n2.Sub(n3))
	as.Float(2.25, n5.Sub(n6))
	as.Integer(15, n1.Sub(n3))
}

func TestInfiniteNumbers(t *testing.T) {
	as := assert.New(t)

	as.False(I(98).IsPosInf())
	as.False(I(0).IsNegInf())

	posInf := F(1).Div(F(0))
	negInf := F(-1).Div(F(0))

	as.True(posInf.IsPosInf())
	as.False(posInf.IsNegInf())
	as.True(negInf.IsNegInf())
	as.False(negInf.IsPosInf())

	as.Compare(data.GreaterThan, posInf, I(1))
	as.Compare(data.LessThan, negInf, I(1))
	as.Compare(data.LessThan, I(1), posInf)
	as.Compare(data.GreaterThan, I(1), negInf)

	as.Compare(data.LessThan, negInf, F(1))
	as.Compare(data.GreaterThan, posInf, F(1))
	as.Compare(data.LessThan, F(1), posInf)
	as.Compare(data.GreaterThan, F(1), negInf)
}

func TestNonNumbers(t *testing.T) {
	as := assert.New(t)

	nan := F(math.Log(-1.0))

	as.True(nan.IsNaN())
	as.False(F(35.5).IsNaN())
	as.False(I(35).IsNaN())

	as.Compare(data.Incomparable, F(1), nan)
	as.Compare(data.Incomparable, nan, F(1))
	as.Compare(data.Incomparable, I(1), nan)
	as.Compare(data.Incomparable, nan, I(1))
}

func TestStringifyNumbers(t *testing.T) {
	as := assert.New(t)
	n1 := data.ParseFloat("12.8")
	n2 := F(12.9)
	n3 := F(20)

	as.String("12.8", n1)
	as.String("12.9", n2)
	as.String("20", n3)
}
