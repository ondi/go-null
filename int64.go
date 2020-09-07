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

// swagger:type integer
type Int64 struct {
	// swagger:ignore
	Int64 int64
	// swagger:ignore
	Valid bool
}

func (self Int64) String() string {
	if self.Valid {
		return strconv.FormatInt(self.Int64, 10)
	}
	return "null"
}

func (self Int64) MarshalJSON() ([]byte, error) {
	if self.Valid {
		return json.Marshal(self.Int64)
	}
	return json.Marshal(nil)
}

func (self *Int64) UnmarshalJSON(data []byte) (err error) {
	var temp *int64
	if err = json.Unmarshal(data, &temp); err != nil {
		return
	}
	if temp != nil {
		self.Int64, self.Valid = *temp, true
	} else {
		self.Int64, self.Valid = 0, false
	}
	return
}

func (self *Int64) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case nil:
		self.Int64, self.Valid = 0, false
		return
	case int64:
		self.Int64, self.Valid = v, true
		return
	case bool:
		if v {
			self.Int64, self.Valid = 1, true
		} else {
			self.Int64, self.Valid = 0, true
		}
		return
	default:
		return fmt.Errorf("not supported: %T %v", v, v)
	}
}

func (self Int64) Value() (driver.Value, error) {
	if self.Valid {
		return self.Int64, nil
	}
	return nil, nil
}
