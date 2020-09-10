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
	test1 := Time{time.Date(2020, 9, 10, 11, 12, 13, 14, time.FixedZone("UTC+3", 3*60*60)), true}
	assert.Assert(t, test1.String() == "2020-09-10T11:12:13+03:00", test1.String())

	test2 := Time{time.Date(2020, 9, 10, 11, 12, 13, 14, time.FixedZone("UTC+3", 3*60*60)), false}
	assert.Assert(t, test2.String() == "null")
}

func TestTimeTs(t *testing.T) {
	test1 := TimeTs{Time{time.Date(2020, 9, 10, 11, 12, 13, 14, time.FixedZone("UTC+3", 3*60*60)), true}}
	assert.Assert(t, test1.String() == "1599725533")

	test2 := TimeTs{Time{time.Date(2020, 9, 10, 11, 12, 13, 14, time.FixedZone("UTC+3", 3*60*60)), false}}
	assert.Assert(t, test2.String() == "null")
}
