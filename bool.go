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

func (self Bool) IsEmptyJSON() bool {
	return self.Valid == false
}

func (self Bool) String(quotes ...string) string {
	if self.Valid {
		if len(quotes) > 1 {
			return quotes[0] + strconv.FormatBool(self.Data) + quotes[1]
		}
		return strconv.FormatBool(self.Data)
	}
	return "null"
}

func (self Bool) StringInt(quotes ...string) (res string) {
	if self.Valid {
		if self.Data {
			res = "1"
		} else {
			res = "0"
		}
		if len(quotes) > 1 {
			return quotes[0] + res + quotes[1]
		}
		return
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
	case nil:
		self.Valid = false
	case bool:
		self.Data, self.Valid = v, true
	case int64:
		if v == 0 {
			self.Data, self.Valid = false, true
		} else {
			self.Data, self.Valid = true, true
		}
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
