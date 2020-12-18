//
//
//

package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

// swagger:type string
type String []string

func Err(err error) (res String) {
	if err != nil {
		res = String{err.Error()}
	}
	return
}

func (self String) Valid() bool {
	return len(self) != 0
}

func (self String) Get() string {
	if len(self) != 0 {
		return self[0]
	}
	return ""
}

func (self String) IsEmptyJSON() bool {
	return len(self) == 0
}

func (self String) String(quotes ...string) string {
	if len(self) != 0 {
		if len(quotes) > 1 {
			return quotes[0] + self[0] + quotes[1]
		}
		return self[0]
	}
	return "null"
}

func (self String) StringSql(quotes ...string) (res string) {
	if len(self) != 0 {
		res = strings.NewReplacer(
			"'", "''",
			"\r", "\\r",
			"\n", "\\n",
		).Replace(self[0])
		if len(quotes) > 1 {
			return quotes[0] + res + quotes[1]
		}
		return
	}
	return "null"
}

func (self String) MarshalJSON() ([]byte, error) {
	if len(self) != 0 {
		return json.Marshal(self[0])
	}
	return json.Marshal(nil)
}

func (self *String) UnmarshalJSON(data []byte) (err error) {
	var temp *string
	if err = json.Unmarshal(data, &temp); err != nil {
		return
	}
	if temp != nil {
		*self = String{*temp}
	} else {
		*self = String{}
	}
	return
}

func (self *String) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var temp *string
	if err = unmarshal(&temp); err != nil {
		return
	}
	if temp != nil {
		*self = String{*temp}
	} else {
		*self = String{}
	}
	return
}

func (self *String) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case nil:
		*self = String{}
		return
	case string:
		*self = String{v}
		return
	case []uint8:
		*self = String{string(v)}
		return
	default:
		return fmt.Errorf("not supported: %T %v", value, value)
	}
}

func (self String) Value() (driver.Value, error) {
	if len(self) != 0 {
		return self[0], nil
	}
	return nil, nil
}
