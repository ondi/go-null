//
//
//

package null

import (
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestString(t *testing.T) {
	test1 := String{"lalala", true}
	assert.Assert(t, test1.String() == "lalala")

	test2 := String{"lalala", false}
	assert.Assert(t, test2.String() == "null")
}

func TestInt64(t *testing.T) {
	test1 := Int64{10, true}
	assert.Assert(t, test1.String() == "10")

	test2 := Int64{10, false}
	assert.Assert(t, test2.String() == "null")
}

func TestFloat64(t *testing.T) {
	test1 := Float64{5.5, true}
	assert.Assert(t, test1.String() == "5.5e+00")

	test2 := Float64{5.5, false}
	assert.Assert(t, test2.String() == "null")
}

func TestBool(t *testing.T) {
	test1 := Bool{false, true}
	assert.Assert(t, test1.String() == "false")

	test2 := Bool{false, false}
	assert.Assert(t, test2.String() == "null")
}

func TestTime(t *testing.T) {
	test1 := Time{time.Time{}, true}
	assert.Assert(t, test1.String() == "0001-01-01T00:00:00Z")

	test2 := Time{time.Time{}, false}
	assert.Assert(t, test2.String() == "null")
}
