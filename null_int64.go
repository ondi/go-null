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

func FormatInt(in int64) string {
	return strconv.FormatInt(in, 10)
}

type Int64 struct {
	Data  int64 `json:"-"`
	Valid bool  `json:"-"`
}

func (self Int64) IsZero() bool {
	return !self.Valid
}

func (self Int64) String() string {
	if self.Valid {
		return FormatInt(self.Data)
	}
	return "null"
}

func (self Int64) Strings(not_valid string, format func(in int64) string, op ...StringOption) (res string) {
	if self.Valid {
		res = format(self.Data)
		for _, v := range op {
			res = v(res)
		}
		return
	}
	return not_valid
}

func (self Int64) MarshalJSON() ([]byte, error) {
	if self.Valid {
		return []byte(FormatInt(self.Data)), nil
	}
	return []byte("null"), nil
}

func (self *Int64) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 || data[0] == 'n' {
		self.Valid = false
		return
	}
	if self.Data, err = strconv.ParseInt(string(data), 0, 64); err == nil {
		self.Valid = true
		return
	}
	err = fmt.Errorf("Int64: %s %w", data, err)
	return
}

func (self *Int64) UnmarshalYAML(value *yaml.Node) (err error) {
	var temp *int64
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

func (self *Int64) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case int64:
		self.Data, self.Valid = v, true
	case string:
		if self.Data, err = strconv.ParseInt(v, 0, 64); err == nil {
			self.Valid = true
		}
	case []uint8:
		if self.Data, err = strconv.ParseInt(string(v), 0, 64); err == nil {
			self.Valid = true
		}
	case float64:
		self.Data, self.Valid = int64(v), true
	case bool:
		if v {
			self.Data, self.Valid = 1, true
		} else {
			self.Data, self.Valid = 0, true
		}
	case nil:
		self.Valid = false
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
