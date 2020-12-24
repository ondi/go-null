//
//
//

package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

var TimeFormatIn = []string{
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02T15:04:05Z0700",
	"2006-01-02T15:04:05",
	"2006-01-02T15:04",
	"2006-01-02",
	"2006-01-02 15:04:05 -07:00",
	"2006-01-02 15:04:05 -0700",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04",
}

var TimeFormatOut = "2006-01-02T15:04:05Z07:00"

// swagger:type string
// swagger:strfmt date-time
type Time []time.Time

func (self Time) Valid() bool {
	return len(self) != 0
}

func (self Time) Get() time.Time {
	if len(self) != 0 {
		return self[0]
	}
	return time.Time{}
}

func (self Time) IsEmptyJSON() bool {
	return len(self) == 0
}

func (self Time) String(quotes ...string) string {
	if len(self) != 0 {
		if len(quotes) > 1 {
			return quotes[0] + self[0].Format(TimeFormatOut) + quotes[1]
		}
		return self[0].Format(TimeFormatOut)
	}
	return "null"
}

func (self Time) MarshalJSON() ([]byte, error) {
	if len(self) != 0 {
		return json.Marshal(self[0].Format(TimeFormatOut))
	}
	return []byte("null"), nil
}

func (self *Time) UnmarshalJSON(data []byte) (err error) {
	if len(data) < 2 {
		err = fmt.Errorf("FORMAT")
		return
	}
	if data[0] == 'n' || data[1] == '"' {
		*self = Time{}
		return
	}
	var temp time.Time
	str := string(data[1 : len(data)-1])
	for _, layout := range TimeFormatIn {
		if temp, err = time.Parse(layout, str); err == nil {
			*self = Time{temp}
			return
		}
	}
	return
}

func (self *Time) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var temp string
	if err = unmarshal(&temp); err != nil {
		return
	}
	if len(temp) == 0 || temp == "null" {
		*self = Time{}
		return
	}
	var test time.Time
	for _, layout := range TimeFormatIn {
		if test, err = time.Parse(layout, temp); err == nil {
			*self = Time{test}
			return
		}
	}
	return
}

func (self *Time) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case nil:
		*self = Time{}
		return
	case time.Time:
		*self = Time{v}
		return
	default:
		return fmt.Errorf("not supported: %T %v", value, value)
	}
}

func (self Time) Value() (driver.Value, error) {
	if len(self) != 0 {
		return self[0], nil
	}
	return nil, nil
}

// swagger:type integer
type TimeUnix struct {
	Time
}

func (self TimeUnix) String() string {
	if len(self.Time) != 0 {
		return strconv.FormatInt(self.Time[0].Unix(), 10)
	}
	return "null"
}

func (self TimeUnix) MarshalJSON() ([]byte, error) {
	if len(self.Time) != 0 {
		return json.Marshal(self.Time[0].Unix())
	}
	return []byte("null"), nil
}

func (self *TimeUnix) UnmarshalJSON(data []byte) (err error) {
	var res int64
	if res, err = strconv.ParseInt(string(data), 0, 64); err != nil {
		return
	}
	self.Time = Time{time.Unix(res, 0)}
	return
}

// swagger:type integer
type TimeUnixNano struct {
	Time
}

func (self TimeUnixNano) String() string {
	if len(self.Time) != 0 {
		return strconv.FormatInt(self.Time[0].UnixNano(), 10)
	}
	return "null"
}

func (self TimeUnixNano) MarshalJSON() ([]byte, error) {
	if len(self.Time) != 0 {
		return json.Marshal(self.Time[0].UnixNano())
	}
	return []byte("null"), nil
}

func (self *TimeUnixNano) UnmarshalJSON(data []byte) (err error) {
	var res int64
	if res, err = strconv.ParseInt(string(data), 0, 64); err != nil {
		return
	}
	self.Time = Time{time.Unix(0, res)}
	return
}
