//
//
//

package null

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strings"
)

type FromString_t struct {
	state         func(r rune) (err error)
	Int           int64
	Exp           int64
	frac_exp      int64
	sign_int      bool
	sign_exp      bool
	frac_overflow bool
}

func (self *FromString_t) parse_int1(r rune) (err error) {
	switch r {
	case '0':
		self.state = self.parse_int2
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		self.Int = int64(r - '0')
		self.state = self.parse_int4
	case '-':
		self.sign_int = true
		self.state = self.parse_int3
	case '+':
		self.state = self.parse_int3
	case '.':
		self.state = self.parse_frac1
	default:
		err = fmt.Errorf("parse_int1: invalid format %q", r)
	}
	return
}

// check format 0x, 0b
func (self *FromString_t) parse_int2(r rune) (err error) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		self.Int = self.Int*10 + int64(r-'0')
		self.state = self.parse_int4
	case 'b', 'B':
		err = fmt.Errorf("binary format not supported")
	case 'x', 'X':
		err = fmt.Errorf("hexadecimal format not supported")
	case '.':
		self.state = self.parse_frac2
	case 0:
		self.state = nil
	default:
		err = fmt.Errorf("parse_int2: invalid format %q", r)
	}
	return
}

// expecting a digit or fraction
func (self *FromString_t) parse_int3(r rune) (err error) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		self.Int = self.Int*10 + int64(r-'0')
		self.state = self.parse_int4
	case '.':
		self.state = self.parse_frac1
	default:
		err = fmt.Errorf("parse_int3: invalid format %q", r)
	}
	return
}

func (self *FromString_t) parse_int4(r rune) (err error) {
	var ok bool
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if self.Int, ok = MulAddInt64(self.Int, 10, int64(r-'0')); !ok {
			err = fmt.Errorf("parse_int4: overflow")
			return
		}
		self.state = self.parse_int4
	case '.':
		self.state = self.parse_frac2
	case 'e', 'E':
		self.state = self.parse_exp1
	case 0:
		self.state = nil
	default:
		err = fmt.Errorf("parse_int4: invalid format %q", r)
	}
	return
}

// expecting a digit
func (self *FromString_t) parse_frac1(r rune) (err error) {
	var ok bool
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if self.Int, ok = MulAddInt64(self.Int, 10, int64(r-'0')); !ok {
			if self.frac_overflow {
				err = fmt.Errorf("parse_frac1: overflow")
			} else {
				self.state = nil
			}
			return
		}
		self.frac_exp++
		self.state = self.parse_frac2
	default:
		err = fmt.Errorf("parse_frac1: invalid format %q", r)
	}
	return
}

func (self *FromString_t) parse_frac2(r rune) (err error) {
	var ok bool
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if self.Int, ok = MulAddInt64(self.Int, 10, int64(r-'0')); !ok {
			if self.frac_overflow {
				err = fmt.Errorf("parse_frac2: overflow")
			} else {
				self.state = nil
			}
			return
		}
		self.frac_exp++
		self.state = self.parse_frac2
	case 'e', 'E':
		self.state = self.parse_exp1
	case 0:
		self.state = nil
	default:
		err = fmt.Errorf("parse_frac2: invalid format %q", r)
	}
	return
}

func (self *FromString_t) parse_exp1(r rune) (err error) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		self.Exp = int64(r - '0')
		self.state = self.parse_exp3
	case '-':
		self.sign_exp = true
		self.state = self.parse_exp2
	case '+':
		self.state = self.parse_exp2
	default:
		err = fmt.Errorf("parse_exp1: invalid format %q", r)
	}
	return
}

// expecting a digit
func (self *FromString_t) parse_exp2(r rune) (err error) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		self.Exp = self.Exp*10 + int64(r-'0')
		self.state = self.parse_exp3
	default:
		err = fmt.Errorf("parse_exp2: invalid format %q", r)
	}
	return
}

func (self *FromString_t) parse_exp3(r rune) (err error) {
	var ok bool
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if self.Exp, ok = MulAddInt64(self.Exp, 10, int64(r-'0')); !ok {
			err = fmt.Errorf("parse_exp3: overflow")
			return
		}
		self.state = self.parse_exp3
	case 0:
		self.state = nil
	default:
		err = fmt.Errorf("parse_exp3: invalid format %q", r)
	}
	return
}

func (self *FromString_t) final() {
	if self.sign_exp {
		self.Exp = -self.Exp
	}
	self.Exp -= self.frac_exp
	if self.sign_int {
		self.Int = -self.Int
	}
}

func ParseFloatString(in string, frac_overflow bool) (Int int64, Exp int64, err error) {
	p, err := ParseFloatReader(strings.NewReader(in), frac_overflow)
	if err != nil {
		err = fmt.Errorf("%q - %w", in, err)
	}
	return p.Int, p.Exp, err
}

func ParseFloatByte(in []byte, frac_overflow bool) (Int int64, Exp int64, err error) {
	p, err := ParseFloatReader(bytes.NewReader(in), frac_overflow)
	if err != nil {
		err = fmt.Errorf("%q - %w", in, err)
	}
	return p.Int, p.Exp, err
}

func ParseFloatReader(reader io.RuneReader, frac_overflow bool) (p FromString_t, err error) {
	p.state = p.parse_int1
	p.frac_overflow = frac_overflow
	for p.state != nil {
		r, _, _ := reader.ReadRune()
		if err = p.state(r); err != nil {
			return
		}
	}
	p.final()
	return
}

// https://wiki.sei.cmu.edu/confluence/display/c/INT32-C.+Ensure+that+operations+on+signed+integers+do+not+result+in+overflow

func AddInt64(a int64, b int64) (int64, bool) {
	if b > 0 {
		if a > math.MaxInt64-b {
			return a, false
		}
	} else if b < 0 {
		if a < math.MinInt64-b {
			return a, false
		}
	}
	return a + b, true
}

func SubInt64(a int64, b int64) (int64, bool) {
	if b > 0 {
		if a < math.MinInt64+b {
			return a, false
		}
	} else if b < 0 {
		if a > math.MaxInt64+b {
			return a, false
		}
	}
	return a - b, true
}

func MulInt64(a int64, b int64) (int64, bool) {
	if a > 0 {
		if b > 0 {
			if a > math.MaxInt64/b {
				return a, false
			}
		} else {
			if b < math.MinInt64/a {
				return a, false
			}
		}
	} else if a < 0 {
		if b > 0 {
			if a < math.MinInt64/b {
				return a, false
			}
		} else {
			if b < math.MaxInt64/a {
				return a, false
			}
		}
	}
	return a * b, true
}

// res = a * b + c
func MulAddInt64(a int64, b int64, c int64) (int64, bool) {
	if temp, ok := MulInt64(a, b); ok {
		if temp, ok = AddInt64(temp, c); ok {
			return temp, true
		}
	}
	return a, false
}
