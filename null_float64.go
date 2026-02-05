//
//
//

package null

import (
	"database/sql/driver"
	"fmt"
	"strconv"

	"gopkg.in/yaml.v3"
)

func FormatFloat(in float64) string {
	return strconv.FormatFloat(in, 'g', -1, 64)
}

type Float64 struct {
	Data  float64 `json:"-"`
	Valid bool    `json:"-"`
}

func (self Float64) IsZero() bool {
	return !self.Valid
}

func (self Float64) String() string {
	if self.Valid {
		return FormatFloat(self.Data)
	}
	return "null"
}

func (self Float64) Strings(not_valid string, format func(in float64) string, op ...StringOption) (res string) {
	if self.Valid {
		res = format(self.Data)
		for _, v := range op {
			res = v(res)
		}
		return
	}
	return not_valid
}

func (self Float64) MarshalJSON() ([]byte, error) {
	if self.Valid {
		return []byte(FormatFloat(self.Data)), nil
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

func (self *Float64) UnmarshalYAML(value *yaml.Node) (err error) {
	var temp *float64
	if err = value.Decode(&temp); err != nil {
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
