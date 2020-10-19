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

// swagger:type string
// swagger:strfmt date-time
type Time struct {
	// swagger:ignore
	Time time.Time
	// swagger:ignore
	Valid bool
}

func (self Time) IsEmptyJSON() bool {
	return !self.Valid
}

func (self Time) String(quotes ...string) string {
	if self.Valid {
		if len(quotes) > 1 {
			return quotes[0] + self.Time.Format("2006-01-02T15:04:05Z07:00") + quotes[1]
		}
		return self.Time.Format("2006-01-02T15:04:05Z07:00")
	}
	return "null"
}

func (self Time) MarshalJSON() ([]byte, error) {
	if self.Valid {
		return json.Marshal(self.Time.Format("2006-01-02T15:04:05Z07:00"))
	}
	return json.Marshal(nil)
}

func (self *Time) UnmarshalJSON(data []byte) (err error) {
	data_str := string(data)
	self.Time, self.Valid = time.Time{}, false
	if len(data_str) == 0 || data_str == "null" || data_str == `""` {
		return
	}
	if self.Time, err = time.Parse(`"2006-01-02T15:04:05Z07:00"`, data_str); err != nil {
		if self.Time, err = time.Parse(`"2006-01-02T15:04:05Z0700"`, data_str); err != nil {
			if self.Time, err = time.Parse(`"2006-01-02T15:04:05"`, data_str); err != nil {
				if self.Time, err = time.Parse(`"2006-01-02T15:04"`, data_str); err != nil {
					if self.Time, err = time.Parse(`"2006-01-02 15:04:05 -07:00"`, data_str); err != nil {
						if self.Time, err = time.Parse(`"2006-01-02 15:04:05 -0700"`, data_str); err != nil {
							if self.Time, err = time.Parse(`"2006-01-02 15:04:05"`, data_str); err != nil {
								if self.Time, err = time.Parse(`"2006-01-02 15:04"`, data_str); err != nil {
									if self.Time, err = time.Parse(`"2006-01-02"`, data_str); err != nil {
										return
									}
								}
							}
						}
					}
				}
			}
		}
	}
	self.Valid = true
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
type TimeTs struct {
	// swagger:ignore
	Time
}

func (self TimeTs) String() string {
	if self.Valid {
		return strconv.FormatInt(self.Time.Time.Unix(), 10)
	}
	return "null"
}

func (self TimeTs) MarshalJSON() ([]byte, error) {
	if self.Valid {
		return json.Marshal(self.Time.Time.Unix())
	}
	return json.Marshal(nil)
}

func (self *TimeTs) UnmarshalJSON(data []byte) (err error) {
	var res int64
	if res, err = strconv.ParseInt(string(data), 0, 64); err != nil {
		return
	}
	self.Time.Time, self.Time.Valid = time.Unix(res, 0), true
	return
}
