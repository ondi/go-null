//
//
//

package null

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
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
	if len(in) <= limit {
		return in
	}
	var r rune
	for ; limit > 0; limit-- {
		if r, _ = utf8.DecodeLastRuneInString(in[:limit]); r != utf8.RuneError {
			break
		}
	}
	return in[:limit]
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
	case string:
		self.Data, self.Valid = v, true
	case []uint8:
		self.Data, self.Valid = string(v), true
	case int64:
		self.Data, self.Valid = strconv.FormatInt(v, 10), true
	case float64:
		self.Data, self.Valid = strconv.FormatFloat(v, 'e', -1, 64), true
	case time.Time:
		self.Data, self.Valid = v.Format(TimeFormatOut), true
	case bool:
		if v {
			self.Data, self.Valid = "true", true
		} else {
			self.Data, self.Valid = "false", true
		}
	case nil:
		self.Valid = false
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
type Str = String

// StringPrice uses no quotes in string representation
// swagger:type string
type StringPrice struct {
	Str
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
