//
//
//

package null

import (
	"encoding/json"
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestString(t *testing.T) {
	var test1 String
	json.Unmarshal([]byte(`"lalala"`), &test1)
	assert.Assert(t, test1.String() == "lalala", test1)

	temp, err := json.Marshal(test1)
	assert.NilError(t, err)
	assert.Assert(t, string(temp) == `"lalala"`)

	json.Unmarshal([]byte(`null`), &test1)
	assert.Assert(t, test1.String() == "null", test1)
}

func TestInt64(t *testing.T) {
	var test1 Int64
	json.Unmarshal([]byte("10"), &test1)
	assert.Assert(t, test1.String() == "10", test1)

	temp, err := json.Marshal(test1)
	assert.NilError(t, err)
	assert.Assert(t, string(temp) == "10", test1)

	json.Unmarshal([]byte("null"), &test1)
	assert.Assert(t, test1.String() == "null", test1)
}

func TestFloat64(t *testing.T) {
	var test1 Float64
	json.Unmarshal([]byte("5.5"), &test1)
	assert.Assert(t, test1.String() == "5.5e+00")

	temp, err := json.Marshal(test1)
	assert.NilError(t, err)
	assert.Assert(t, string(temp) == "5.5e+00", test1)

	json.Unmarshal([]byte("null"), &test1)
	assert.Assert(t, test1.String() == "null")
}

func TestBool(t *testing.T) {
	var test1 Bool
	json.Unmarshal([]byte("false"), &test1)
	assert.Assert(t, test1.String() == "false")

	temp, err := json.Marshal(test1)
	assert.NilError(t, err)
	assert.Assert(t, string(temp) == "false", test1)

	json.Unmarshal([]byte("null"), &test1)
	assert.Assert(t, test1.String() == "null")
}

func TestTime01(t *testing.T) {
	in := time.Date(2020, 9, 10, 11, 12, 13, 14, time.FixedZone("UTC+3", 3*60*60))
	test1 := Time{Data: in, Valid: true}
	assert.Assert(t, test1.String() == "2020-09-10T11:12:13+03:00", test1.String())
	assert.Assert(t, test1.Data == in)

	test2 := Time{}
	assert.Assert(t, test2.String() == "null")

	var err error
	var test3 Time

	err = json.Unmarshal([]byte("\"2020-09-10T11:12:13+03:00\""), &test3)
	assert.NilError(t, err)
	assert.Assert(t, test3.String() == "2020-09-10T11:12:13+03:00", test3.String())

	err = json.Unmarshal([]byte("\"2020-09-10T11:12:13+0300\""), &test3)
	assert.NilError(t, err)
	assert.Assert(t, test3.String() == "2020-09-10T11:12:13+03:00", test3.String())

	err = json.Unmarshal([]byte("\"2020-09-10T11:12:13\""), &test3)
	assert.NilError(t, err)
	assert.Assert(t, test3.String() == "2020-09-10T11:12:13Z", test3.String())

	err = json.Unmarshal([]byte("\"2020-09-10T11:12\""), &test3)
	assert.NilError(t, err)
	assert.Assert(t, test3.String() == "2020-09-10T11:12:00Z", test3.String())

	err = json.Unmarshal([]byte("\"2020-09-10 11:12:13 +03:00\""), &test3)
	assert.NilError(t, err)
	assert.Assert(t, test3.String() == "2020-09-10T11:12:13+03:00", test3.String())

	err = json.Unmarshal([]byte("\"2020-09-10 11:12:13 +0300\""), &test3)
	assert.NilError(t, err)
	assert.Assert(t, test3.String() == "2020-09-10T11:12:13+03:00", test3.String())

	err = json.Unmarshal([]byte("\"2020-09-10 11:12:13\""), &test3)
	assert.NilError(t, err)
	assert.Assert(t, test3.String() == "2020-09-10T11:12:13Z", test3.String())

	err = json.Unmarshal([]byte("\"2020-09-10 11:12\""), &test3)
	assert.NilError(t, err)
	assert.Assert(t, test3.String() == "2020-09-10T11:12:00Z", test3.String())

	err = json.Unmarshal([]byte("\"2020-09-10\""), &test3)
	assert.NilError(t, err)
	assert.Assert(t, test3.String() == "2020-09-10T00:00:00Z", test3.String())

	err = json.Unmarshal([]byte("\"2020-09\""), &test3)
	assert.Assert(t, err != nil, "should be error")
}

func TestTime02(t *testing.T) {
	var test1 struct {
		Test Time `json:"test"`
	}
	in1 := "null"
	err := json.Unmarshal([]byte(in1), &test1)
	assert.NilError(t, err)
	assert.Assert(t, test1.Test.String() == in1, test1.Test.String())

	var test2 struct {
		Test Time `json:"test"`
	}
	in2 := `{"test":null}`
	err2 := json.Unmarshal([]byte(in2), &test2)
	assert.NilError(t, err2)
	assert.Assert(t, test2.Test.String() == "null", test2.Test.String())

	var test3 struct {
		Test Time `json:"test"`
	}
	in3 := `{"test":""}`
	err3 := json.Unmarshal([]byte(in3), &test3)
	assert.NilError(t, err3)
	assert.Assert(t, test3.Test.String() == "null", test3.Test.String())

	var test4 struct {
		Test Time `json:"test"`
	}
	in4 := `{"test":"2020-01-02T15:04:05+03:00"}`
	err4 := json.Unmarshal([]byte(in4), &test4)
	assert.NilError(t, err4)
	assert.Assert(t, test4.Test.String() == "2020-01-02T15:04:05+03:00", test4.Test.String())
}

func TestTimeUnix(t *testing.T) {
	in := time.Date(2020, 9, 10, 11, 12, 13, 14, time.FixedZone("UTC+3", 3*60*60))
	test1 := TimeUnix{Time{Data: in, Valid: true}}
	assert.Assert(t, test1.String() == "1599725533")
	assert.Assert(t, test1.Time.Data == in)

	test2 := TimeUnix{}
	assert.Assert(t, test2.String() == "null")
}

func TestTimeUnixNano(t *testing.T) {
	in := time.Date(2020, 9, 10, 11, 12, 13, 14, time.FixedZone("UTC+3", 3*60*60))
	test1 := TimeUnixNano{Time{Data: in, Valid: true}}
	assert.Assert(t, test1.String() == "1599725533000000014")
	assert.Assert(t, test1.Data == in)

	test2 := TimeUnixNano{}
	assert.Assert(t, test2.String() == "null")
}
