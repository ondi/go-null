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

// swagger:type number
type Float64 struct {
	// swagger:ignore
	Float64 float64
	// swagger:ignore
	Valid bool
}

func (self Float64) String() string {
	if self.Valid {
		return strconv.FormatFloat(self.Float64, 'e', -1, 64)
	}
	return "null"
}

func (self Float64) MarshalJSON() ([]byte, error) {
	if self.Valid {
		return json.Marshal(self.Float64)
	}
	return json.Marshal(nil)
}

func (self *Float64) UnmarshalJSON(data []byte) (err error) {
	var temp *float64
	if err = json.Unmarshal(data, &temp); err != nil {
		return
	}
	if temp != nil {
		self.Float64, self.Valid = *temp, true
	} else {
		self.Float64, self.Valid = 0, false
	}
	return
}

func (self *Float64) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case nil:
		self.Float64, self.Valid = 0, false
		return
	case float64:
		self.Float64, self.Valid = v, true
		return
	case []uint8:
		if self.Float64, err = strconv.ParseFloat(string(v), 64); err == nil {
			self.Valid = true
		}
		return
	default:
		return fmt.Errorf("not supported: %T %v", value, value)
	}
}

func (self Float64) Value() (driver.Value, error) {
	if self.Valid {
		return self.Float64, nil
	}
	return nil, nil
}
