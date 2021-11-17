//
//
//

package null

import (
	"database/sql/driver"
)

// swagger:type boolean
type None struct{}

func (self None) String() string {
	return "null"
}

func (self None) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}

func (self *None) UnmarshalJSON(data []byte) (err error) {
	return
}

func (self *None) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	return
}

func (self *None) Scan(value interface{}) (err error) {
	return
}

func (self None) Value() (driver.Value, error) {
	return nil, nil
}
