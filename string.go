//
//
//

package null

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

var Replacer = strings.NewReplacer(
	"'", "''",
)

// swagger:type string
type String struct {
	// swagger:ignore
	Data string
	// swagger:ignore
	Valid bool
}

func Err(err error) (res String) {
	if err != nil {
		res.Data, res.Valid = err.Error(), true
	}
	return
}

func StringLimit(in string, limit int) string {
	if len(in) > limit {
		var prev int
		for i := range in {
			if i > limit {
				return in[:prev]
			}
			prev = i
		}
	}
	return in
}

func (self String) String() string {
	if self.Valid {
		return self.Data
	}
	return "null"
}

func (self String) StringQuote(a string, b string) string {
	if self.Valid {
		return a + self.Data + b
	}
	return "null"
}

func (self String) StringSql(a string, b string) string {
	if self.Valid {
		return a + Replacer.Replace(self.Data) + b
	}
	return "null"
}

func (self String) StringSqlLimit(a string, b string, limit int) (res string) {
	if self.Valid {
		return a + Replacer.Replace(StringLimit(self.Data, limit)) + b
	}
	return "null"
}

func (self String) MarshalJSON() ([]byte, error) {
	if self.Valid {
		return []byte(strconv.Quote(self.Data)), nil
	}
	return []byte("null"), nil
}

func (self *String) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 || data[0] == 'n' {
		self.Valid = false
		return
	}
	if self.Data, err = strconv.Unquote(string(data)); err != nil {
		err = fmt.Errorf("%.32s: %w", data, err)
		return
	}
	self.Valid = true
	return
}

func (self *String) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var temp *string
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

func (self *String) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case nil:
		self.Valid = false
	case string:
		self.Data, self.Valid = v, true
	case []uint8:
		self.Data, self.Valid = string(v), true
	default:
		err = fmt.Errorf("not supported: %T %v", value, value)
	}
	return
}

func (self String) Value() (driver.Value, error) {
	if self.Valid {
		return self.Data, nil
	}
	return nil, nil
}

// allow to use .String() method with embedded String struct
type str = String

// StringPrice uses no quotes in string representation
// swagger:type string
type StringPrice struct {
	str
}

func (self StringPrice) MarshalJSON() (res []byte, err error) {
	if self.Valid {
		temp := strconv.Quote(self.Data)
		return []byte(temp[1 : len(temp)-1]), nil
	}
	return []byte("null"), nil
}

func (self *StringPrice) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 || data[0] == 'n' {
		self.Valid = false
		return
	}
	if self.Data, err = strconv.Unquote(`"` + string(data) + `"`); err != nil {
		err = fmt.Errorf("%.32s: %w", data, err)
		return
	}
	self.Valid = true
	return
}
