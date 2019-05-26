package builtin_test

import (
	"testing"
	"time"

	"gitlab.com/kode4food/ale/core/builtin"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
)

func TestCurrentTime(t *testing.T) {
	as := assert.New(t)

	t1 := time.Now().UnixNano()
	t2 := int64(builtin.CurrentTime().(data.Integer))

	as.Equal(t1-(t1%1000000), t2-(t2%1000000))
}
