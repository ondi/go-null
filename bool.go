//
//
//

package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

// swagger:type boolean
type Bool struct {
	// swagger:ignore
	Bool bool
	// swagger:ignore
	Valid bool
}

func (self Bool) String() string {
	if self.Valid {
		return strconv.FormatBool(self.Bool)
	}
	return "null"
}

func (self Bool) MarshalJSON() ([]byte, error) {
	if self.Valid {
		return json.Marshal(self.Bool)
	}
	return json.Marshal(nil)
}

func (self *Bool) UnmarshalJSON(data []byte) (err error) {
	var temp *bool
	if err = json.Unmarshal(data, &temp); err != nil {
		return
	}
	if temp != nil {
		self.Bool, self.Valid = *temp, true
	} else {
		self.Bool, self.Valid = false, false
	}
	return
}

func (self *Bool) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case nil:
		self.Bool, self.Valid = false, false
		return
	case bool:
		self.Bool, self.Valid = v, true
		return
	case int64:
		if v == 0 {
			self.Bool, self.Valid = false, true
		} else {
			self.Bool, self.Valid = true, true
		}
		return
	default:
		return fmt.Errorf("not supported: %T %v", value, value)
	}
}

func (self Bool) Value() (driver.Value, error) {
	if self.Valid {
		return self.Bool, nil
	}
	return nil, nil
}
