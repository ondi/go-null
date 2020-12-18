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
type Bool []bool

func (self Bool) Valid() bool {
	return len(self) != 0
}

func (self Bool) Get() bool {
	if len(self) != 0 {
		return self[0]
	}
	return false
}

func (self Bool) IsEmptyJSON() bool {
	return len(self) == 0
}

func (self Bool) String(quotes ...string) string {
	if len(self) != 0 {
		if len(quotes) > 1 {
			return quotes[0] + strconv.FormatBool(self[0]) + quotes[1]
		}
		return strconv.FormatBool(self[0])
	}
	return "null"
}

func (self Bool) StringInt(quotes ...string) (res string) {
	if len(self) != 0 {
		if self[0] {
			res = "1"
		} else {
			res = "0"
		}
		if len(quotes) > 1 {
			return quotes[0] + res + quotes[1]
		}
		return
	}
	return "null"
}

func (self Bool) MarshalJSON() ([]byte, error) {
	if len(self) != 0 {
		return json.Marshal(self[0])
	}
	return json.Marshal(nil)
}

func (self *Bool) UnmarshalJSON(data []byte) (err error) {
	var temp *bool
	if err = json.Unmarshal(data, &temp); err != nil {
		return
	}
	if temp != nil {
		*self = Bool{*temp}
	} else {
		*self = Bool{}
	}
	return
}

func (self *Bool) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var temp *bool
	if err = unmarshal(&temp); err != nil {
		return
	}
	if temp != nil {
		*self = Bool{*temp}
	} else {
		*self = Bool{}
	}
	return
}

func (self *Bool) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case nil:
		*self = Bool{}
		return
	case bool:
		*self = Bool{v}
		return
	case int64:
		if v == 0 {
			*self = Bool{false}
		} else {
			*self = Bool{true}
		}
		return
	default:
		return fmt.Errorf("not supported: %T %v", value, value)
	}
}

func (self Bool) Value() (driver.Value, error) {
	if len(self) != 0 {
		return self[0], nil
	}
	return nil, nil
}
