//
//
//

package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// swagger:type string
// swagger:strfmt string
type String struct {
	// swagger:ignore
	Str string
	// swagger:ignore
	Valid bool
}

func (self String) String() string {
	if self.Valid {
		return self.Str
	}
	return "null"
}

func (self String) MarshalJSON() ([]byte, error) {
	if self.Valid {
		return json.Marshal(self.Str)
	}
	return json.Marshal(nil)
}

func (self *String) UnmarshalJSON(data []byte) (err error) {
	var temp *string
	if err = json.Unmarshal(data, &temp); err != nil {
		return
	}
	if temp != nil {
		self.Str, self.Valid = *temp, true
	} else {
		self.Str, self.Valid = "", false
	}
	return
}

func (self *String) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case nil:
		self.Str, self.Valid = "", false
		return
	case string:
		self.Str, self.Valid = v, true
		return
	case []uint8:
		self.Str, self.Valid = string(v), true
		return
	default:
		return fmt.Errorf("not supported: %T %v", value, value)
	}
}

func (self String) Value() (driver.Value, error) {
	if self.Valid {
		return self.Str, nil
	}
	return nil, nil
}
