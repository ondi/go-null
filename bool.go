//
//
//

package null

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

type Bool struct {
	Data  bool `json:"-"`
	Valid bool `json:"-"`
}

func FormatBool(in bool) string {
	return strconv.FormatBool(in)
}

func (self Bool) String() string {
	if self.Valid {
		return FormatBool(self.Data)
	}
	return "null"
}

func (self Bool) Strings(not_valid string, format func(in bool) string) string {
	if self.Valid {
		return format(self.Data)
	}
	return not_valid
}

func (self Bool) MarshalJSON() ([]byte, error) {
	if self.Valid {
		return []byte(strconv.FormatBool(self.Data)), nil
	}
	return []byte("null"), nil
}

func (self *Bool) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 || data[0] == 'n' {
		self.Valid = false
		return
	}
	if self.Data, err = strconv.ParseBool(string(data)); err == nil {
		self.Valid = true
		return
	}
	err = fmt.Errorf("Bool: %s %w", data, err)
	return
}

func (self *Bool) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var temp *bool
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

func (self *Bool) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case bool:
		self.Data, self.Valid = v, true
	case int64:
		if v == 0 {
			self.Data, self.Valid = false, true
		} else {
			self.Data, self.Valid = true, true
		}
	case string:
		if self.Data, err = strconv.ParseBool(v); err == nil {
			self.Valid = true
		}
	case []uint8:
		if self.Data, err = strconv.ParseBool(string(v)); err == nil {
			self.Valid = true
		}
	case nil:
		self.Valid = false
	default:
		err = fmt.Errorf("not supported: %T %v", value, value)
	}
	return
}

func (self Bool) Value() (driver.Value, error) {
	if self.Valid {
		return self.Data, nil
	}
	return nil, nil
}
