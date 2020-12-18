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
type Float64 []float64

func (self Float64) Valid() bool {
	return len(self) != 0
}

func (self Float64) Get() float64 {
	if len(self) != 0 {
		return self[0]
	}
	return 0
}

func (self Float64) IsEmptyJSON() bool {
	return len(self) == 0
}

func (self Float64) String(quotes ...string) string {
	if len(self) != 0 {
		if len(quotes) > 1 {
			return quotes[0] + strconv.FormatFloat(self[0], 'e', -1, 64) + quotes[1]
		}
		return strconv.FormatFloat(self[0], 'e', -1, 64)
	}
	return "null"
}

func (self Float64) MarshalJSON() ([]byte, error) {
	if len(self) != 0 {
		return json.Marshal(self[0])
	}
	return json.Marshal(nil)
}

func (self *Float64) UnmarshalJSON(data []byte) (err error) {
	var temp *float64
	if err = json.Unmarshal(data, &temp); err != nil {
		return
	}
	if temp != nil {
		*self = Float64{*temp}
	} else {
		*self = Float64{}
	}
	return
}

func (self *Float64) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var temp *float64
	if err = unmarshal(&temp); err != nil {
		return
	}
	if temp != nil {
		*self = Float64{*temp}
	} else {
		*self = Float64{}
	}
	return
}

func (self *Float64) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case nil:
		*self = Float64{}
		return
	case float64:
		*self = Float64{v}
		return
	case []uint8:
		var temp float64
		if temp, err = strconv.ParseFloat(string(v), 64); err == nil {
			*self = Float64{temp}
		}
		return
	default:
		return fmt.Errorf("not supported: %T %v", value, value)
	}
}

func (self Float64) Value() (driver.Value, error) {
	if len(self) != 0 {
		return self[0], nil
	}
	return nil, nil
}
