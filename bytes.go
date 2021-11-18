//
//
//

package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

var ByteReplace = func() ByteOption {
	return func(in []byte) []byte {
		return bytes.ReplaceAll(in, []byte{'\''}, []byte{'\'', '\''})
	}
}

// swagger:type string
type Bytes struct {
	// swagger:ignore
	Data []byte
	// swagger:ignore
	Valid bool
}

func (self Bytes) String() string {
	if self.Valid {
		return string(self.Data)
	}
	return "null"
}

func (self Bytes) Bytes(op ...ByteOption) []byte {
	if self.Valid {
		for _, v := range op {
			self.Data = v(self.Data)
		}
		return self.Data
	}
	return []byte("null")
}

type ByteOption func(in []byte) []byte

func ByteLimit(limit int) ByteOption {
	return func(in []byte) []byte {
		if len(in) > limit {
			return in[:limit]
		}
		return in
	}
}

func ByteQuote(a string, b string) ByteOption {
	return func(in []byte) []byte {
		var temp bytes.Buffer
		temp.WriteString(a)
		temp.Write(in)
		temp.WriteString(b)
		return temp.Bytes()
	}
}

func (self Bytes) MarshalJSON() ([]byte, error) {
	if self.Valid {
		return json.Marshal(self.Data)
	}
	return []byte("null"), nil
}

func (self *Bytes) UnmarshalJSON(data []byte) (err error) {
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

func (self *Bytes) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var temp *[]byte
	if err = unmarshal(&temp); err != nil {
		return
	}
	if temp != nil {
		self.Data, self.Valid = *temp, true
	} else {
		self.Valid = false
	}
	return
}

func (self *Bytes) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case string:
		self.Data, self.Valid = []byte(v), true
	case []uint8:
		self.Data, self.Valid = append([]byte{}, v...), true
	case nil:
		self.Valid = false
	default:
		err = fmt.Errorf("not supported: %T %v", value, value)
	}
	return
}

func (self Bytes) Value() (driver.Value, error) {
	if self.Valid {
		return self.Data, nil
	}
	return nil, nil
}
