//
//
//

package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type String struct {
	Data  string `json:"-"`
	Valid bool   `json:"-"`
}

func StringLimit(in string, limit int) string {
	if len(in) > limit {
		for ; limit > 0; limit-- {
			if r, _ := utf8.DecodeLastRuneInString(in[:limit]); r != utf8.RuneError {
				break
			}
		}
		return in[:limit]
	}
	return in
}

func (self String) String() string {
	if self.Valid {
		return self.Data
	}
	return "null"
}

func (self String) Strings(not_valid string, op ...StringOption) string {
	if self.Valid {
		for _, v := range op {
			self.Data = v(self.Data)
		}
		return self.Data
	}
	return not_valid
}

type StringOption func(in string) string

func StrLimit(limit int) StringOption {
	return func(in string) string {
		return StringLimit(in, limit)
	}
}

func StrEscape() StringOption {
	return func(in string) string {
		return strings.NewReplacer("'", "''", "\x00", "\\x00", "\x1a", "\\x1a", "\\", "\\\\").Replace(in)
	}
}

func StrSqlQuote() StringOption {
	return func(in string) string {
		return "'" + StrEscape()(in) + "'"
	}
}

func (self String) MarshalJSON() ([]byte, error) {
	if self.Valid {
		return json.Marshal(self.Data)
	}
	return []byte("null"), nil
}

func (self *String) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 || data[0] == 'n' {
		self.Valid = false
		return
	}
	if err = json.Unmarshal(data, &self.Data); err == nil {
		self.Valid = true
		return
	}
	err = fmt.Errorf("String: %s %w", data, err)
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

// make .String() method available for embedded String struct
type Str = String

// StringPrice uses no quotes in string representation
type StringPrice struct {
	Str
}

func (self StringPrice) MarshalJSON() (res []byte, err error) {
	if self.Valid {
		return []byte(self.Data), nil
	}
	return []byte("null"), nil
}

func (self *StringPrice) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 || data[0] == 'n' {
		self.Valid = false
		return
	}
	if data[0] == '"' && json.Unmarshal(data, &self.Data) == nil {
		self.Valid = true
	} else {
		self.Data, self.Valid = string(data), true
	}
	return
}
