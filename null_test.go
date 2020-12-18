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
	test1 := String{"lalala"}
	assert.Assert(t, test1.String() == "lalala")

	test2 := String{}
	assert.Assert(t, test2.String() == "null")
}

func TestInt64(t *testing.T) {
	test1 := Int64{10}
	assert.Assert(t, test1.String() == "10")

	test2 := Int64{}
	assert.Assert(t, test2.String() == "null")
}

func TestFloat64(t *testing.T) {
	test1 := Float64{5.5}
	assert.Assert(t, test1.String() == "5.5e+00")

	test2 := Float64{}
	assert.Assert(t, test2.String() == "null")
}

func TestBool(t *testing.T) {
	test1 := Bool{false}
	assert.Assert(t, test1.String() == "false")

	test2 := Bool{}
	assert.Assert(t, test2.String() == "null")
}

func TestTime(t *testing.T) {
	test1 := Time{time.Date(2020, 9, 10, 11, 12, 13, 14, time.FixedZone("UTC+3", 3*60*60))}
	assert.Assert(t, test1.String() == "2020-09-10T11:12:13+03:00", test1.String())

	test2 := Time{}
	assert.Assert(t, test2.String() == "null")

	var err error
	var test3 Time

	err = test3.UnmarshalJSON([]byte("\"2020-09-10T11:12:13+03:00\""))
	assert.NilError(t, err)
	assert.Assert(t, test3.String() == "2020-09-10T11:12:13+03:00", test3.String())

	err = test3.UnmarshalJSON([]byte("\"2020-09-10T11:12:13+0300\""))
	assert.NilError(t, err)
	assert.Assert(t, test3.String() == "2020-09-10T11:12:13+03:00", test3.String())

	err = test3.UnmarshalJSON([]byte("\"2020-09-10T11:12:13\""))
	assert.NilError(t, err)
	assert.Assert(t, test3.String() == "2020-09-10T11:12:13Z", test3.String())

	err = test3.UnmarshalJSON([]byte("\"2020-09-10T11:12\""))
	assert.NilError(t, err)
	assert.Assert(t, test3.String() == "2020-09-10T11:12:00Z", test3.String())

	err = test3.UnmarshalJSON([]byte("\"2020-09-10 11:12:13 +03:00\""))
	assert.NilError(t, err)
	assert.Assert(t, test3.String() == "2020-09-10T11:12:13+03:00", test3.String())

	err = test3.UnmarshalJSON([]byte("\"2020-09-10 11:12:13 +0300\""))
	assert.NilError(t, err)
	assert.Assert(t, test3.String() == "2020-09-10T11:12:13+03:00", test3.String())

	err = test3.UnmarshalJSON([]byte("\"2020-09-10 11:12:13\""))
	assert.NilError(t, err)
	assert.Assert(t, test3.String() == "2020-09-10T11:12:13Z", test3.String())

	err = test3.UnmarshalJSON([]byte("\"2020-09-10 11:12\""))
	assert.NilError(t, err)
	assert.Assert(t, test3.String() == "2020-09-10T11:12:00Z", test3.String())

	err = test3.UnmarshalJSON([]byte("\"2020-09-10\""))
	assert.NilError(t, err)
	assert.Assert(t, test3.String() == "2020-09-10T00:00:00Z", test3.String())

	err = test3.UnmarshalJSON([]byte("\"2020-09\""))
	assert.Assert(t, err != nil, "should be error")
}

func TestTimeUnix(t *testing.T) {
	test1 := TimeUnix{time.Date(2020, 9, 10, 11, 12, 13, 14, time.FixedZone("UTC+3", 3*60*60))}
	assert.Assert(t, test1.String() == "1599725533")

	test2 := TimeUnix{}
	assert.Assert(t, test2.String() == "null")
}

func TestTimeUnixNano(t *testing.T) {
	test1 := TimeUnixNano{time.Date(2020, 9, 10, 11, 12, 13, 14, time.FixedZone("UTC+3", 3*60*60))}
	assert.Assert(t, test1.String() == "1599725533000000014")

	test2 := TimeUnixNano{}
	assert.Assert(t, test2.String() == "null")
}
