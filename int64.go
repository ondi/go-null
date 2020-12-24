//
//
//

package null

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

// swagger:type integer
type Int64 struct {
	// swagger:ignore
	Data int64
	// swagger:ignore
	Valid bool
}

func (self Int64) IsEmptyJSON() bool {
	return self.Valid == false
}

func (self Int64) String(quotes ...string) string {
	if self.Valid {
		if len(quotes) > 1 {
			return quotes[0] + strconv.FormatInt(self.Data, 10) + quotes[1]
		}
		return strconv.FormatInt(self.Data, 10)
	}
	return "null"
}

func (self Int64) MarshalJSON() ([]byte, error) {
	if self.Valid {
		return []byte(strconv.FormatInt(self.Data, 10)), nil
	}
	return []byte("null"), nil
}

func (self *Int64) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 || data[0] == 'n' {
		self.Valid = false
		return
	}
	if self.Data, err = strconv.ParseInt(string(data), 0, 64); err != nil {
		return
	}
	self.Valid = true
	return
}

func (self *Int64) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var temp *int64
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

func (self *Int64) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case nil:
		self.Valid = false
	case int64:
		self.Data, self.Valid = v, true
	case bool:
		if v {
			self.Data, self.Valid = 1, true
		} else {
			self.Valid = false
		}
	default:
		err = fmt.Errorf("not supported: %T %v", v, v)
	}
	return
}

func (self Int64) Value() (driver.Value, error) {
	if self.Valid {
		return self.Data, nil
	}
	return nil, nil
}
