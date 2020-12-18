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
type Int64 []int64

func (self Int64) Valid() bool {
	return len(self) != 0
}

func (self Int64) Get() int64 {
	if len(self) != 0 {
		return self[0]
	}
	return 0
}

func (self Int64) IsEmptyJSON() bool {
	return len(self) == 0
}

func (self Int64) String(quotes ...string) string {
	if len(self) != 0 {
		if len(quotes) > 1 {
			return quotes[0] + strconv.FormatInt(self[0], 10) + quotes[1]
		}
		return strconv.FormatInt(self[0], 10)
	}
	return "null"
}

func (self Int64) MarshalJSON() ([]byte, error) {
	if len(self) != 0 {
		return json.Marshal(self[0])
	}
	return json.Marshal(nil)
}

func (self *Int64) UnmarshalJSON(data []byte) (err error) {
	var temp *int64
	if err = json.Unmarshal(data, &temp); err != nil {
		return
	}
	if temp != nil {
		*self = Int64{*temp}
	} else {
		*self = Int64{}
	}
	return
}

func (self *Int64) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var temp *int64
	if err = unmarshal(&temp); err != nil {
		return
	}
	if temp != nil {
		*self = Int64{*temp}
	} else {
		*self = Int64{}
	}
	return
}

func (self *Int64) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case nil:
		*self = Int64{}
		return
	case int64:
		*self = Int64{v}
		return
	case bool:
		if v {
			*self = Int64{1}
		} else {
			*self = Int64{0}
		}
		return
	default:
		return fmt.Errorf("not supported: %T %v", v, v)
	}
}

func (self Int64) Value() (driver.Value, error) {
	if len(self) != 0 {
		return self[0], nil
	}
	return nil, nil
}
