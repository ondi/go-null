//
// behaves like string
//

package null

import (
	"database/sql/driver"
	"fmt"

	"gopkg.in/yaml.v3"
)

var FRAC_OVERFLOW = false

type Decimal64 struct {
	Int   int64 `json:"-"`
	Exp   int64 `json:"-"`
	Valid bool  `json:"-"`
}

func (self Decimal64) IsZero() bool {
	return !self.Valid
}

func (self Decimal64) String() string {
	if self.Valid {
		return fmt.Sprintf("%de%d", self.Int, self.Exp)
	}
	return "null"
}

func (self Decimal64) Strings(not_valid string, op ...StringOption) string {
	temp := self.String()
	if self.Valid {
		for _, v := range op {
			temp = v(temp)
		}
		return temp
	}
	return not_valid
}

func (self Decimal64) MarshalJSON() ([]byte, error) {
	if self.Valid {
		return []byte(self.String()), nil
	}
	return []byte("null"), nil
}

func (self *Decimal64) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 || data[0] == 'n' {
		self.Valid = false
		return
	}
	if data[0] == '"' {
		self.Int, self.Exp, err = ParseFloatByte(data[1:len(data)-1], FRAC_OVERFLOW)
	} else {
		self.Int, self.Exp, err = ParseFloatByte(data, FRAC_OVERFLOW)
	}
	if err == nil {
		self.Valid = true
	}
	return
}

func (self *Decimal64) UnmarshalYAML(value *yaml.Node) (err error) {
	var temp *string
	if err = value.Decode(&temp); err != nil {
		return
	}
	if temp != nil {
		if self.Int, self.Exp, err = ParseFloatString(*temp, FRAC_OVERFLOW); err == nil {
			self.Valid = true
		}
	} else {
		self.Valid = false
	}
	return
}

func (self *Decimal64) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case string:
		if self.Int, self.Exp, err = ParseFloatString(v, FRAC_OVERFLOW); err == nil {
			self.Valid = true
		}
	case []uint8:
		if self.Int, self.Exp, err = ParseFloatByte(v, FRAC_OVERFLOW); err == nil {
			self.Valid = true
		}
	case int64:
		self.Int, self.Exp, self.Valid = v, 0, true
	case float64:
		if self.Int, self.Exp, err = ParseFloatFloat(v, FRAC_OVERFLOW); err == nil {
			self.Valid = true
		}
	case bool:
		if v {
			self.Int, self.Exp, self.Valid = 1, 0, true
		} else {
			self.Int, self.Exp, self.Valid = 0, 0, true
		}
	case nil:
		self.Valid = false
	default:
		err = fmt.Errorf("not supported: %T %v", value, value)
	}
	return
}

func (self Decimal64) Value() (driver.Value, error) {
	if self.Valid {
		return self.String(), nil
	}
	return nil, nil
}

func (self *Decimal64) IntPart(exp_bias int64) (res int64, ok bool) {
	ok = true
	res = self.Int
	exp := self.Exp + exp_bias
	if exp < 0 {
		for i := exp; i < 0; i++ {
			res = res / 10
		}
	} else if exp > 0 {
		for i := int64(0); i < exp; i++ {
			if res, ok = MulInt64(res, 10); !ok {
				return
			}
		}
	}
	return
}
