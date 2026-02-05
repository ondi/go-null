//
//
//

package null

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

var TimeFormatIn = []string{
	`2006-01-02T15:04:05.999999999Z07:00`,
	`2006-01-02T15:04:05.999999999Z0700`,
	`2006-01-02T15:04:05.999999999Z07`,
	`2006-01-02T15:04:05.999999999Z07:00:00`,
	`2006-01-02T15:04:05.999999999`,
	`2006-01-02T15:04`,
	`2006-01-02`,
	`2006-01-02 15:04:05.999999999 -07:00`,
	`2006-01-02 15:04:05.999999999 -0700`,
	`2006-01-02 15:04:05.999999999 -07`,
	`2006-01-02 15:04:05.999999999 -07:00:00`,
	`2006-01-02 15:04:05.999999999-07:00`,
	`2006-01-02 15:04:05.999999999-0700`,
	`2006-01-02 15:04:05.999999999-07`,
	`2006-01-02 15:04:05.999999999-07:00:00`,
	`2006-01-02 15:04:05.999999999`,
	`2006-01-02 15:04`,
}

// time.RFC3339Nano
var TimeFormatOut = "2006-01-02T15:04:05.999999999Z07:00"

func FormatTime(in time.Time) string {
	return in.Format(TimeFormatOut)
}

type Time struct {
	Data  time.Time `json:"-"`
	Valid bool      `json:"-"`
}

func (self Time) IsZero() bool {
	return !self.Valid
}

func (self Time) String() string {
	if self.Valid {
		return FormatTime(self.Data)
	}
	return "null"
}

func (self Time) Strings(not_valid string, format func(time.Time) string, op ...StringOption) (res string) {
	if self.Valid {
		res = format(self.Data)
		for _, v := range op {
			res = v(res)
		}
		return
	}
	return not_valid
}

func (self Time) MarshalJSON() ([]byte, error) {
	if self.Valid {
		return []byte(`"` + FormatTime(self.Data) + `"`), nil
	}
	return []byte("null"), nil
}

func (self *Time) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 || data[0] == 'n' || (len(data) == 2 && data[0] == '"' && data[1] == '"') {
		self.Valid = false
		return
	}
	str := string(data)
	if len(str) > 2 {
		str = str[1 : len(str)-1]
	}
	for _, layout := range TimeFormatIn {
		if self.Data, err = time.Parse(layout, str); err == nil {
			self.Valid = true
			return
		}
	}
	err = fmt.Errorf("Time: %s %w", data, err)
	return
}

func (self *Time) UnmarshalYAML(value *yaml.Node) (err error) {
	var temp *string
	if err = value.Decode(&temp); err != nil {
		return
	}
	if temp == nil {
		self.Valid = false
		return
	}
	for _, layout := range TimeFormatIn {
		if self.Data, err = time.Parse(layout, *temp); err == nil {
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
	case []uint8:
		str := string(v)
		for _, layout := range TimeFormatIn {
			if self.Data, err = time.Parse(layout, str); err == nil {
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

type TimeUnix struct {
	Time
}

func (self TimeUnix) String() string {
	if self.Valid {
		return FormatInt(self.Data.Unix())
	}
	return "null"
}

func (self TimeUnix) MarshalJSON() ([]byte, error) {
	if self.Valid {
		return []byte(FormatInt(self.Data.Unix())), nil
	}
	return []byte("null"), nil
}

func (self *TimeUnix) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 || data[0] == 'n' {
		self.Valid = false
		return
	}
	res, err := strconv.ParseInt(string(data), 0, 64)
	if err != nil {
		return
	}
	self.Data, self.Valid = time.Unix(res, 0), true
	return
}

type TimeUnixNano struct {
	Time
}

func (self TimeUnixNano) String() string {
	if self.Valid {
		return FormatInt(self.Data.UnixNano())
	}
	return "null"
}

func (self TimeUnixNano) MarshalJSON() ([]byte, error) {
	if self.Valid {
		return []byte(FormatInt(self.Data.UnixNano())), nil
	}
	return []byte("null"), nil
}

func (self *TimeUnixNano) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 || data[0] == 'n' {
		self.Valid = false
		return
	}
	res, err := strconv.ParseInt(string(data), 0, 64)
	if err != nil {
		return
	}
	self.Data, self.Valid = time.Unix(0, res), true
	return
}
