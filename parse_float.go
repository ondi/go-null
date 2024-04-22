//
//
//

package null

import (
	"fmt"
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
const (
	sign = 0b_10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000
	expn = 0b_01111111_11110000_00000000_00000000_00000000_00000000_00000000_00000000
	frac = 0b_00000000_00001111_11111111_11111111_11111111_11111111_11111111_11111111
)

type state_t func(r rune, size int) (state state_t, err error)

type Decimal_t struct {
	Int      int64
	Exp      int64
	frac_exp int64
	sign_int bool
	sign_exp bool
}

func (self *Decimal_t) parse_int1(r rune, size int) (state state_t, err error) {
	switch r {
	case '0':
		state = self.parse_int2
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		self.Int = int64(r - '0')
		state = self.parse_int4
	case '-':
		self.sign_int = true
		state = self.parse_int3
	case '+':
		state = self.parse_int3
	case '.':
		state = self.parse_frac1
	default:
		err = fmt.Errorf("parse_int1: invalid format %q", r)
	}
	return
}

// check format 0x, 0b
func (self *Decimal_t) parse_int2(r rune, size int) (state state_t, err error) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		self.Int = self.Int*10 + int64(r-'0')
		state = self.parse_int4
	case 'b', 'B':
		err = fmt.Errorf("binary format not supported")
	case 'x', 'X':
		err = fmt.Errorf("hexadecimal format not supported")
	case '.':
		state = self.parse_frac1
	case 0:
		// ok
	default:
		err = fmt.Errorf("parse_int2: invalid format %q", r)
	}
	return
}

// here should not be EOF
func (self *Decimal_t) parse_int3(r rune, size int) (state state_t, err error) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		self.Int = self.Int*10 + int64(r-'0')
		state = self.parse_int4
	case '.':
		state = self.parse_frac1
	default:
		err = fmt.Errorf("parse_int2: invalid format %q", r)
	}
	return
}

func (self *Decimal_t) parse_int4(r rune, size int) (state state_t, err error) {
	var ok bool
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if self.Int, ok = MulAdd64(self.Int, 10, int64(r-'0')); !ok {
			err = fmt.Errorf("parse_int3: overflow")
			return
		}
		state = self.parse_int4
	case '.':
		state = self.parse_frac1
	case 0:
		// ok
	default:
		err = fmt.Errorf("parse_int3: invalid format %q", r)
	}
	return
}

func (self *Decimal_t) parse_frac1(r rune, size int) (state state_t, err error) {
	var ok bool
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if self.Int, ok = MulAdd64(self.Int, 10, int64(r-'0')); !ok {
			err = fmt.Errorf("parse_frac1: overflow")
			return
		}
		self.frac_exp++
		state = self.parse_frac1
	case 'e', 'E':
		state = self.parse_exp1
	case 0:
		// ok
	default:
		err = fmt.Errorf("parse_frac1: invalid format %q", r)
	}
	return
}

func (self *Decimal_t) parse_exp1(r rune, size int) (state state_t, err error) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		self.Exp = int64(r - '0')
		state = self.parse_exp3
	case '-':
		self.sign_exp = true
		state = self.parse_exp2
	case '+':
		state = self.parse_exp2
	default:
		err = fmt.Errorf("parse_exp1: invalid format %q", r)
	}
	return
}

// here should not be EOF
func (self *Decimal_t) parse_exp2(r rune, size int) (state state_t, err error) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		self.Exp = self.Exp*10 + int64(r-'0')
		state = self.parse_exp3
	default:
		err = fmt.Errorf("parse_exp2: invalid format %q", r)
	}
	return
}

func (self *Decimal_t) parse_exp3(r rune, size int) (state state_t, err error) {
	var ok bool
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if self.Exp, ok = MulAdd64(self.Exp, 10, int64(r-'0')); !ok {
			err = fmt.Errorf("parse_exp3: overflow")
			return
		}
		state = self.parse_exp3
	case 0:
		// ok
	default:
		err = fmt.Errorf("parse_exp3: invalid format %q", r)
	}
	return
}

func (self *Decimal_t) final() {
	if self.sign_exp {
		self.Exp = -self.Exp
	}
	self.Exp -= self.frac_exp
	if self.sign_int {
		self.Int = -self.Int
	}
}

func (self *Decimal_t) String() string {
	if self.Exp == 0 {
		return fmt.Sprintf("%d", self.Int)
	}
	return fmt.Sprintf("%de%d", self.Int, self.Exp)
}

func (self *Decimal_t) Int64() (int64, bool) {
	if self.Exp < 0 {
		return self.Int / Width10(-self.Exp), true
	}
	return Mul64(self.Int, Width10(self.Exp))
}

func ParseFloat(in string) (res Decimal_t, err error) {
	state := res.parse_int1
	reader := strings.NewReader(in)
	for state != nil {
		last_rune, last_size, _ := reader.ReadRune()
		if state, err = state(last_rune, last_size); err != nil {
			err = fmt.Errorf("%v - %w", in, err)
			return
		}
	}
	res.final()
	return
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
