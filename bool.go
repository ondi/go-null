//
//
//

package null

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

// swagger:type boolean
type Bool struct {
	// swagger:ignore
	Data bool
	// swagger:ignore
	Valid bool
}

func (self Bool) String() string {
	if self.Valid {
		return strconv.FormatBool(self.Data)
	}
	return "null"
}

func (self Bool) StringInt() string {
	if self.Valid {
		if self.Data {
			return "1"
		}
		return "0"
	}
	return "null"
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
	if self.Data, err = strconv.ParseBool(string(data)); err != nil {
		return
	}
	self.Valid = true
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
