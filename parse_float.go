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

// IEEE 754
// Double: SEEEEEEE EEEEFFFF FFFFFFFF FFFFFFFF FFFFFFFF FFFFFFFF FFFFFFFF FFFFFFFF
// 1 [63]	11 [62–52]	52 [51–00]
// The exponent field is an 11-bit unsigned integer from 0 to 2047, in biased form:
// an exponent value of 1023 represents the actual zero.
// Exponents range from −1022 to +1023 because exponents
// of −1023 (all 0s) and +1024 (all 1s) are reserved for special numbers.
// const (
// 	sign = 0b_10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000
// 	expn = 0b_01111111_11110000_00000000_00000000_00000000_00000000_00000000_00000000
// 	frac = 0b_00000000_00001111_11111111_11111111_11111111_11111111_11111111_11111111
// )

type ParseFloat_t struct {
	Int           int64
	Exp           int64
	state         func(r rune, size int) (err error)
	frac_exp      int64
	sign_int      bool
	sign_exp      bool
	frac_overflow bool
}

func (self *ParseFloat_t) parse_int1(r rune, size int) (err error) {
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
func (self *ParseFloat_t) parse_int2(r rune, size int) (err error) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		self.Int = self.Int*10 + int64(r-'0')
		self.state = self.parse_int4
	case 'b', 'B':
		err = fmt.Errorf("binary format not supported")
	case 'x', 'X':
		err = fmt.Errorf("hexadecimal format not supported")
	case '.':
		self.state = self.parse_frac1
	case 0:
		self.state = nil
	default:
		err = fmt.Errorf("parse_int2: invalid format %q", r)
	}
	return
}

// here should not be EOF
func (self *ParseFloat_t) parse_int3(r rune, size int) (err error) {
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

func (self *ParseFloat_t) parse_int4(r rune, size int) (err error) {
	var ok bool
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if self.Int, ok = MulAdd64(self.Int, 10, int64(r-'0')); !ok {
			err = fmt.Errorf("parse_int4: overflow")
			return
		}
		self.state = self.parse_int4
	case '.':
		self.state = self.parse_frac1
	case 0:
		self.state = nil
	default:
		err = fmt.Errorf("parse_int4: invalid format %q", r)
	}
	return
}

func (self *ParseFloat_t) parse_frac1(r rune, size int) (err error) {
	var ok bool
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if self.Int, ok = MulAdd64(self.Int, 10, int64(r-'0')); !ok {
			if self.frac_overflow {
				self.state = nil
			} else {
				err = fmt.Errorf("parse_frac1: overflow")
			}
			return
		}
		self.frac_exp++
		self.state = self.parse_frac1
	case 'e', 'E':
		self.state = self.parse_exp1
	case 0:
		self.state = nil
	default:
		err = fmt.Errorf("parse_frac1: invalid format %q", r)
	}
	return
}

func (self *ParseFloat_t) parse_exp1(r rune, size int) (err error) {
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

// here should not be EOF
func (self *ParseFloat_t) parse_exp2(r rune, size int) (err error) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		self.Exp = self.Exp*10 + int64(r-'0')
		self.state = self.parse_exp3
	default:
		err = fmt.Errorf("parse_exp2: invalid format %q", r)
	}
	return
}

func (self *ParseFloat_t) parse_exp3(r rune, size int) (err error) {
	var ok bool
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if self.Exp, ok = MulAdd64(self.Exp, 10, int64(r-'0')); !ok {
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

func (self *ParseFloat_t) final() {
	if self.sign_exp {
		self.Exp = -self.Exp
	}
	self.Exp -= self.frac_exp
	if self.sign_int {
		self.Int = -self.Int
	}
}

func ParseFloat(in string, frac_overflow bool) (Int int64, Exp int64, err error) {
	if Int, Exp, err = ParseFloatReader(strings.NewReader(in), frac_overflow); err != nil {
		err = fmt.Errorf("%q - %w", in, err)
	}
	return
}

func ParseFloatByte(in []byte, frac_overflow bool) (Int int64, Exp int64, err error) {
	if Int, Exp, err = ParseFloatReader(bytes.NewReader(in), frac_overflow); err != nil {
		err = fmt.Errorf("%q - %w", in, err)
	}
	return
}

func ParseFloatReader(reader io.RuneReader, frac_overflow bool) (Int int64, Exp int64, err error) {
	p := ParseFloat_t{}
	p.state = p.parse_int1
	p.frac_overflow = frac_overflow
	for p.state != nil {
		last_rune, last_size, _ := reader.ReadRune()
		if err = p.state(last_rune, last_size); err != nil {
			return 0, 0, err
		}
	}
	p.final()
	return p.Int, p.Exp, err
}

func Width10(in int64) (res int64) {
	res = 1
	for i := int64(0); i < in; i++ {
		res *= 10
	}
	return
}

// https://wiki.sei.cmu.edu/confluence/display/c/INT32-C.+Ensure+that+operations+on+signed+integers+do+not+result+in+overflow

func Add64(a int64, b int64) (int64, bool) {
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

func Sub64(a int64, b int64) (int64, bool) {
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

func Mul64(a int64, b int64) (int64, bool) {
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
func MulAdd64(a int64, b int64, c int64) (int64, bool) {
	if temp, ok := Mul64(a, b); ok {
		if temp, ok = Add64(temp, c); ok {
			return temp, true
		}
	}
	return a, false
}
