//
//
//

package null

import (
	"database/sql/driver"
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
type Time struct {
	// swagger:ignore
	Data time.Time
	// swagger:ignore
	Valid bool
}

func (self Time) String() string {
	if self.Valid {
		return self.Data.Format(TimeFormatOut)
	}
	return "null"
}

func (self Time) StringQuote(a string, b string) string {
	if self.Valid {
		return a + self.Data.Format(TimeFormatOut) + b
	}
	return "null"
}

func (self Time) StringQuoteFormat(a string, b string, format string) string {
	if self.Valid {
		return a + self.Data.Format(format) + b
	}
	return "null"
}

func (self Time) MarshalJSON() ([]byte, error) {
	if self.Valid {
		return []byte(`"` + self.Data.Format(TimeFormatOut) + `"`), nil
	}
	return []byte("null"), nil
}

func (self *Time) UnmarshalJSON(data []byte) (err error) {
	if len(data) < 2 {
		return fmt.Errorf("FORMAT")
	}
	if data[0] == 'n' || data[1] == '"' {
		*self = Time{}
		return
	}
	str := string(data[1 : len(data)-1])
	for _, layout := range TimeFormatIn {
		if self.Data, err = time.Parse(layout, str); err == nil {
			self.Valid = true
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
	for _, layout := range TimeFormatIn {
		if self.Data, err = time.Parse(layout, temp); err == nil {
			self.Valid = true
			return
		}
	}
	return
}

func (self *Time) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case time.Time:
		self.Data, self.Valid = v, true
	case string:
		for _, layout := range TimeFormatIn {
			if self.Data, err = time.Parse(layout, v); err == nil {
				self.Valid = true
				return
			}
		}
	case nil:
		self.Valid = false
	default:
		err = fmt.Errorf("not supported: %T %v", value, value)
	}
	return
}

func (self Time) Value() (driver.Value, error) {
	if self.Valid {
		return self.Data, nil
	}
	return nil, nil
}

// swagger:type integer
type TimeUnix struct {
	// swagger:ignore
	Time
}

func (self TimeUnix) String() string {
	if self.Valid {
		return strconv.FormatInt(self.Data.Unix(), 10)
	}
	return "null"
}

func (self TimeUnix) MarshalJSON() ([]byte, error) {
	if self.Valid {
		return []byte(strconv.FormatInt(self.Data.Unix(), 10)), nil
	}
	return []byte("null"), nil
}

func (self *TimeUnix) UnmarshalJSON(data []byte) (err error) {
	var res int64
	if res, err = strconv.ParseInt(string(data), 0, 64); err != nil {
		return
	}
	self.Data, self.Valid = time.Unix(res, 0), true
	return
}

// swagger:type integer
type TimeUnixNano struct {
	// swagger:ignore
	Time
}

func (self TimeUnixNano) String() string {
	if self.Valid {
		return strconv.FormatInt(self.Data.UnixNano(), 10)
	}
	return "null"
}

func (self TimeUnixNano) MarshalJSON() ([]byte, error) {
	if self.Valid {
		return []byte(strconv.FormatInt(self.Data.UnixNano(), 10)), nil
	}
	return []byte("null"), nil
}

func (self *TimeUnixNano) UnmarshalJSON(data []byte) (err error) {
	var res int64
	if res, err = strconv.ParseInt(string(data), 0, 64); err != nil {
		return
	}
	self.Data, self.Valid = time.Unix(0, res), true
	return
}
