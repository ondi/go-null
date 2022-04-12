//
//
//

package null

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

type Float64 struct {
	Data  float64 `json:"-"`
	Valid bool    `json:"-"`
}

func (self Float64) String() string {
	if self.Valid {
		return strconv.FormatFloat(self.Data, 'e', -1, 64)
	}
	return "null"
}

func (self Float64) Strings(not_valid string, format func(in float64) string) string {
	if self.Valid {
		return format(self.Data)
	}
	return not_valid
}

func (self Float64) MarshalJSON() ([]byte, error) {
	if self.Valid {
		return []byte(strconv.FormatFloat(self.Data, 'e', -1, 64)), nil
	}
	return []byte("null"), nil
}

func (self *Float64) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 || data[0] == 'n' {
		self.Valid = false
		return
	}
	if self.Data, err = strconv.ParseFloat(string(data), 64); err == nil {
		self.Valid = true
		return
	}
	err = fmt.Errorf("Float: %s %w", data, err)
	return
}

func (self *Float64) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var temp *float64
	if err = unmarshal(&temp); err != nil {
		return
	}
	if temp != nil {
		self.Data, self.Valid = *temp, true
	} else {
		self.Valid = false
	}
	return
}

func (self *Float64) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case float64:
		self.Data, self.Valid = v, true
	case string:
		if self.Data, err = strconv.ParseFloat(v, 64); err == nil {
			self.Valid = true
		}
	case []uint8:
		if self.Data, err = strconv.ParseFloat(string(v), 64); err == nil {
			self.Valid = true
		}
	case int64:
		self.Data, self.Valid = float64(v), true
	case nil:
		self.Valid = false
	default:
		err = fmt.Errorf("not supported: %T %v", value, value)
	}
	return
}

func (self Float64) Value() (driver.Value, error) {
	if self.Valid {
		return self.Data, nil
	}
	return nil, nil
}
