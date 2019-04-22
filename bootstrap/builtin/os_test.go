package builtin_test

import (
	"testing"
	"time"

	"gitlab.com/kode4food/ale/bootstrap/builtin"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
)

func TestCurrentTime(t *testing.T) {
	as := assert.New(t)

	t1 := time.Now().UnixNano()
	t2 := int64(builtin.CurrentTime().(data.Integer))

	as.Equal(t1-(t1%100000), t2-(t2%100000))
}
