package api_test

import (
	"fmt"
	"testing"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestParseNumber(t *testing.T) {
	as := assert.New(t)
	n1 := api.ParseFloat("12.8")
	n2 := F(12.8)
	n3 := F(12.8)

	as.Equal(n1, n2)
	as.Equal(n2, n3)

	defer as.ExpectPanic(fmt.Sprintf(api.ExpectedFloat, S(`'splosion!`)))
	api.ParseFloat("'splosion!")
}

func TestCompareNumbers(t *testing.T) {
	as := assert.New(t)
	n1 := api.ParseFloat("12.8")
	n2 := F(12.9)
	n3 := F(20)

	as.Compare(api.LessThan, n1, n2)
	as.Compare(api.LessThan, n1, n3)
}

func TestStringifyNumbers(t *testing.T) {
	as := assert.New(t)
	n1 := api.ParseFloat("12.8")
	n2 := F(12.9)
	n3 := F(20)

	as.String("12.8", n1)
	as.String("12.9", n2)
	as.String("20", n3)
}
