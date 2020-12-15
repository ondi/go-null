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

var TimeFormat = []string{
	`"2006-01-02T15:04:05Z07:00"`,
	`"2006-01-02T15:04:05Z0700"`,
	`"2006-01-02T15:04:05"`,
	`"2006-01-02T15:04"`,
	`"2006-01-02"`,
	`"2006-01-02 15:04:05 -07:00"`,
	`"2006-01-02 15:04:05 -0700"`,
	`"2006-01-02 15:04:05"`,
	`"2006-01-02 15:04"`,
}

// swagger:type string
// swagger:strfmt date-time
type Time struct {
	// swagger:ignore
	Time time.Time
	// swagger:ignore
	Valid bool
}

func (self *Time) IsEmptyJSON() bool {
	return self.Valid == false
}

func (self *Time) String(quotes ...string) string {
	if self.Valid {
		if len(quotes) > 1 {
			return quotes[0] + self.Time.Format("2006-01-02T15:04:05Z07:00") + quotes[1]
		}
		return self.Time.Format("2006-01-02T15:04:05Z07:00")
	}
	return "null"
}

func (self *Time) MarshalJSON() ([]byte, error) {
	if self.Valid {
		return json.Marshal(self.Time.Format("2006-01-02T15:04:05Z07:00"))
	}
	return json.Marshal(nil)
}

func (self *Time) UnmarshalJSON(data []byte) (err error) {
	str := string(data)
	if len(str) == 0 || str == "null" || str == `""` {
		return
	}
	for _, layout := range TimeFormat {
		if self.Time, err = time.Parse(layout, str); err == nil {
			self.Valid = true
			return
		}
	}
	return
}

func (self *Time) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var temp *time.Time
	if err = unmarshal(&temp); err != nil {
		return
	}
	if temp != nil {
		self.Time, self.Valid = *temp, true
	} else {
		self.Time, self.Valid = time.Time{}, false
	}
	return
}

func (self *Time) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case nil:
		self.Time, self.Valid = time.Time{}, false
		return
	case time.Time:
		self.Time, self.Valid = v, true
		return
	default:
		return fmt.Errorf("not supported: %T %v", value, value)
	}
}

func (self Time) Value() (driver.Value, error) {
	if self.Valid {
		return self.Time, nil
	}
	return nil, nil
}

// swagger:type integer
type TimeUnix struct {
	// swagger:ignore
	Time
}

func (self *TimeUnix) String() string {
	if self.Valid {
		return strconv.FormatInt(self.Time.Time.Unix(), 10)
	}
	return "null"
}

func (self *TimeUnix) MarshalJSON() ([]byte, error) {
	if self.Valid {
		return json.Marshal(self.Time.Time.Unix())
	}
	return json.Marshal(nil)
}

func (self *TimeUnix) UnmarshalJSON(data []byte) (err error) {
	var res int64
	if res, err = strconv.ParseInt(string(data), 0, 64); err != nil {
		return
	}
	self.Time.Time, self.Time.Valid = time.Unix(res, 0), true
	return
}

type TimeUnixNano struct {
	// swagger:ignore
	Time
}

func (self *TimeUnixNano) String() string {
	if self.Valid {
		return strconv.FormatInt(self.Time.Time.UnixNano(), 10)
	}
	return "null"
}

func (self *TimeUnixNano) MarshalJSON() ([]byte, error) {
	if self.Valid {
		return json.Marshal(self.Time.Time.UnixNano())
	}
	return json.Marshal(nil)
}

func (self *TimeUnixNano) UnmarshalJSON(data []byte) (err error) {
	var res int64
	if res, err = strconv.ParseInt(string(data), 0, 64); err != nil {
		return
	}
	self.Time.Time, self.Time.Valid = time.Unix(0, res), true
	return
}
