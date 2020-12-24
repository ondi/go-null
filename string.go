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
	"\r", "\\r",
	"\n", "\\n",
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

func (self String) IsEmptyJSON() bool {
	return self.Valid == false
}

func (self String) String(quotes ...string) string {
	if self.Valid {
		if len(quotes) > 1 {
			return quotes[0] + self.Data + quotes[1]
		}
		return self.Data
	}
	return "null"
}

func (self String) StringSql(quotes ...string) (res string) {
	if self.Valid {
		Replacer.Replace(self.Data)
		if len(quotes) > 1 {
			return quotes[0] + res + quotes[1]
		}
		return
	}
	return "null"
}

func (self String) MarshalJSON() ([]byte, error) {
	if self.Valid {
		return []byte(strconv.Quote(string(self.Data))), nil
	}
	return []byte("null"), nil
}

func (self *String) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 || data[0] == 'n' {
		self.Valid = false
		return
	}
	if self.Data, err = strconv.Unquote(string(data)); err != nil {
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
