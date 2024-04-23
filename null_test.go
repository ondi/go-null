//
//
//

package null

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"gotest.tools/assert"
)

type TestString_t struct {
	Field1 String `json:"field1,omitempty"`
}

func TestString01(t *testing.T) {
	var test1 String
	json.Unmarshal([]byte(`"lalala"`), &test1)
	assert.Assert(t, test1.String() == "lalala", test1)

	temp, err := json.Marshal(test1)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, string(temp) == `"lalala"`)

	json.Unmarshal([]byte(`null`), &test1)
	assert.Assert(t, test1.String() == "null", test1)
}

func TestString02(t *testing.T) {
	var test2 TestString_t

	val := reflect.ValueOf(test2)
	typ := val.Type()
	assert.Assert(t, reflect.Type(typ).Kind() == reflect.Struct)
	assert.Assert(t, reflect.Zero(typ).Interface() == val.Interface())

	// temp, err := json.Marshal(test2)
	// assert.Assert(t, err==nil, err)
	// assert.Assert(t, string(temp) == `{}`, string(temp))
}

func TestString03(t *testing.T) {
	var err error
	var test1 Decimal64

	err = json.Unmarshal([]byte(`""`), &test1)
	assert.Assert(t, err != nil, err)

	err = json.Unmarshal([]byte(`"123.45"`), &test1)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, test1.String() == "12345e-2", test1)

	err = json.Unmarshal([]byte(`123.45`), &test1)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, test1.String() == "12345e-2", test1)

	temp, err := json.Marshal(test1)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, string(temp) == `12345e-2`, test1)

	json.Unmarshal([]byte(`null`), &test1)
	assert.Assert(t, test1.String() == "null", test1)
}

func TestString04(t *testing.T) {
	test1 := String{Data: `123"456`, Valid: true}
	assert.Assert(t, test1.Strings("null", StrQuote1) == `'123"456'`, test1.Strings("null", StrQuote1))

	test1.Valid = false
	assert.Assert(t, test1.Strings("null", StrQuote1) == "null", test1)
}

func TestString05(t *testing.T) {
	test1 := String{Data: `1234567890`, Valid: true}
	assert.Assert(t, test1.Strings("null", StrLimit(5)) == `12345`, test1.Strings("null", StrLimit(5)))

	test1.Valid = false
	assert.Assert(t, test1.Strings("null", StrLimit(5)) == "null", test1.Strings("null", StrLimit(5)))
}

func TestInt64(t *testing.T) {
	var test1 Int64
	json.Unmarshal([]byte("10"), &test1)
	assert.Assert(t, test1.String() == "10", test1)

	temp, err := json.Marshal(test1)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, string(temp) == "10", test1)

	json.Unmarshal([]byte("null"), &test1)
	assert.Assert(t, test1.String() == "null", test1)
}

func TestFloat64(t *testing.T) {
	var test1 Float64
	json.Unmarshal([]byte("5.5"), &test1)
	assert.Assert(t, test1.String() == "5.5")

	temp, err := json.Marshal(test1)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, string(temp) == "5.5", test1)

	json.Unmarshal([]byte("null"), &test1)
	assert.Assert(t, test1.String() == "null")
}

func TestBool(t *testing.T) {
	var test1 Bool
	json.Unmarshal([]byte("false"), &test1)
	assert.Assert(t, test1.String() == "false")

	temp, err := json.Marshal(test1)
	assert.Assert(t, err == nil, err)
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
	assert.Assert(t, err == nil, err)
	assert.Assert(t, test3.String() == "2020-09-10T11:12:13+03:00", test3.String())

	err = json.Unmarshal([]byte("\"2020-09-10T11:12:13+0300\""), &test3)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, test3.String() == "2020-09-10T11:12:13+03:00", test3.String())

	err = json.Unmarshal([]byte("\"2020-09-10T11:12:13\""), &test3)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, test3.String() == "2020-09-10T11:12:13Z", test3.String())

	err = json.Unmarshal([]byte("\"2020-09-10T11:12\""), &test3)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, test3.String() == "2020-09-10T11:12:00Z", test3.String())

	err = json.Unmarshal([]byte("\"2020-09-10 11:12:13 +03:00\""), &test3)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, test3.String() == "2020-09-10T11:12:13+03:00", test3.String())

	err = json.Unmarshal([]byte("\"2020-09-10 11:12:13 +0300\""), &test3)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, test3.String() == "2020-09-10T11:12:13+03:00", test3.String())

	err = json.Unmarshal([]byte("\"2020-09-10 11:12:13\""), &test3)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, test3.String() == "2020-09-10T11:12:13Z", test3.String())

	err = json.Unmarshal([]byte("\"2020-09-10 11:12\""), &test3)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, test3.String() == "2020-09-10T11:12:00Z", test3.String())

	err = json.Unmarshal([]byte("\"2020-09-10\""), &test3)
	assert.Assert(t, err == nil, err)
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
	assert.Assert(t, err == nil, err)
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

func TestScan01(t *testing.T) {
	var a Int64
	query := map[string][]string{"test": {"123"}}
	ScanQuery(&a, "test", query)
	assert.Assert(t, a.Valid == true)
	assert.Assert(t, a.Data == 123)
}

func TestLimit01(t *testing.T) {
	a := StrLimit(8)("你好嗎")
	assert.Assert(t, a == "你好", a)
}

func TestParseFloat01(t *testing.T) {
	var err error
	var Int, Exp int64

	Int, Exp, err = ParseFloatString("", false)
	assert.Assert(t, err != nil, fmt.Sprintf("int=%v, exp=%v", Int, Exp))

	Int, Exp, err = ParseFloatString("-", false)
	assert.Assert(t, err != nil, fmt.Sprintf("int=%v, exp=%v", Int, Exp))

	Int, Exp, err = ParseFloatString("+", false)
	assert.Assert(t, err != nil, fmt.Sprintf("int=%v, exp=%v", Int, Exp))

	Int, Exp, err = ParseFloatString("0", false)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, Int == 0 && Exp == 0, fmt.Sprintf("int=%v, exp=%v", Int, Exp))

	Int, Exp, err = ParseFloatString("00", false)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, Int == 0 && Exp == 0, fmt.Sprintf("int=%v, exp=%v", Int, Exp))

	Int, Exp, err = ParseFloatString("001", false)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, Int == 1 && Exp == 0, fmt.Sprintf("int=%v, exp=%v", Int, Exp))

	Int, Exp, err = ParseFloatString("-0", false)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, Int == 0 && Exp == 0, fmt.Sprintf("int=%v, exp=%v", Int, Exp))

	Int, Exp, err = ParseFloatString("+0", false)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, Int == 0 && Exp == 0, fmt.Sprintf("int=%v, exp=%v", Int, Exp))

	Int, Exp, err = ParseFloatString("-1", false)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, Int == -1 && Exp == 0, fmt.Sprintf("int=%v, exp=%v", Int, Exp))

	Int, Exp, err = ParseFloatString("+1", false)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, Int == 1 && Exp == 0, fmt.Sprintf("int=%v, exp=%v", Int, Exp))

	Int, Exp, err = ParseFloatString("-1.", false)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, Int == -1 && Exp == 0, fmt.Sprintf("int=%v, exp=%v", Int, Exp))

	Int, Exp, err = ParseFloatString("+1.", false)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, Int == 1 && Exp == 0, fmt.Sprintf("int=%v, exp=%v", Int, Exp))

	Int, Exp, err = ParseFloatString("-1.e", false)
	assert.Assert(t, err != nil, fmt.Sprintf("int=%v, exp=%v", Int, Exp))

	Int, Exp, err = ParseFloatString("-1.e0", false)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, Int == -1 && Exp == 0, fmt.Sprintf("int=%v, exp=%v", Int, Exp))

	Int, Exp, err = ParseFloatString("-123.000", false)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, Int == -123000 && Exp == -3, fmt.Sprintf("int=%v, exp=%v", Int, Exp))

	Int, Exp, err = ParseFloatString("-123.0001", false)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, Int == -1230001 && Exp == -4, fmt.Sprintf("int=%v, exp=%v", Int, Exp))

	Int, Exp, err = ParseFloatString("-123.0001e1", false)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, Int == -1230001 && Exp == -3, fmt.Sprintf("int=%v, exp=%v", Int, Exp))

	Int, Exp, err = ParseFloatString("-123.0001e-1", false)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, Int == -1230001 && Exp == -5, fmt.Sprintf("int=%v, exp=%v", Int, Exp))

	Int, Exp, err = ParseFloatString("-123.0001e10", false)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, Int == -1230001 && Exp == 6, fmt.Sprintf("int=%v, exp=%v", Int, Exp))

	Int, Exp, err = ParseFloatString("-123.0001e-10", false)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, Int == -1230001 && Exp == -14, fmt.Sprintf("int=%v, exp=%v", Int, Exp))

	Int, Exp, err = ParseFloatString("9223372036854775807", false)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, Int == 9223372036854775807 && Exp == 0, fmt.Sprintf("int=%v, exp=%v", Int, Exp))

	Int, Exp, err = ParseFloatString("922337203685477580.7", false)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, Int == 9223372036854775807 && Exp == -1, fmt.Sprintf("int=%v, exp=%v", Int, Exp))

	Int, Exp, err = ParseFloatString("922337203685477580.8", false)
	assert.Assert(t, err != nil, fmt.Sprintf("int=%v, exp=%v", Int, Exp))

	Int, Exp, err = ParseFloatString("922337203685477580.8", true)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, Int == 922337203685477580 && Exp == 0, fmt.Sprintf("int=%v, exp=%v", Int, Exp))

	Int, Exp, err = ParseFloatString("9223372036854775808", false)
	assert.Assert(t, err != nil, fmt.Sprintf("int=%v, exp=%v", Int, Exp))

	Int, Exp, err = ParseFloatString("3.1415926535", false)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, Int == 31415926535 && Exp == -10, fmt.Sprintf("int=%v, exp=%v", Int, Exp))

	d := Decimal64{Int: Int, Exp: Exp, Valid: true}
	i, ok := d.Int64()
	t.Logf("int64=%v, ok=%v, str=%v", i, ok, d.String())
}

func TestParseFloat02(t *testing.T) {
	var err error
	var Int, Exp int64

	Int, Exp, err = ParseFloatFloat(3.1415926535, false)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, Int == 31415926535 && Exp == -10, fmt.Sprintf("int=%v, exp=%v", Int, Exp))
}
