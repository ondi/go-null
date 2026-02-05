//
//
//

package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

type String struct {
	Data  string `json:"-"`
	Valid bool   `json:"-"`
}

func (self String) IsZero() bool {
	return !self.Valid
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

func (self *String) UnmarshalYAML(value *yaml.Node) (err error) {
	var temp *string
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

func (self *String) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case string:
		self.Data, self.Valid = v, true
	case []uint8:
		self.Data, self.Valid = string(v), true
	case int64:
		self.Data, self.Valid = FormatInt(v), true
	case float64:
		self.Data, self.Valid = FormatFloat(v), true
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
