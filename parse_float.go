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

type next_state_t func(r rune, size int) (next_state next_state_t)

type Float_t struct {
	Error    string
	input    string
	Int      int64
	frac     int64
	frac_mul int64
	Exp      int64
	int_sign bool
	exp_sign bool
}

func (self *Float_t) parse_int1(r rune, size int) (next_state next_state_t) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		self.Int = int64(r - '0')
		next_state = self.parse_int3
	case '-':
		self.int_sign = true
		next_state = self.parse_int2
	case '+':
		next_state = self.parse_int2
	case '.':
		next_state = self.parse_frac1
	default:
		self.Error = fmt.Sprintf("parse_int1: invalid format %q", self.input)
	}
	return
}

// here should not be EOF
func (self *Float_t) parse_int2(r rune, size int) (next_state next_state_t) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		self.Int = self.Int*10 + int64(r-'0')
		next_state = self.parse_int3
	case '.':
		next_state = self.parse_frac1
	default:
		self.Error = fmt.Sprintf("parse_int2: invalid format %q", self.input)
	}
	return
}

func (self *Float_t) parse_int3(r rune, size int) (next_state next_state_t) {
	var ok bool
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if self.Int, ok = Mul64(self.Int, 10); !ok {
			self.Error = fmt.Sprintf("parse_int3: overflow %q", self.input)
			return
		}
		if self.Int, ok = Add64(self.Int, int64(r-'0')); !ok {
			self.Error = fmt.Sprintf("parse_int3: overflow %q", self.input)
			return
		}
		next_state = self.parse_int3
	case '.':
		next_state = self.parse_frac1
	case 0:
		// ok
	default:
		self.Error = fmt.Sprintf("parse_int3: invalid format %q", self.input)
	}
	return
}

func (self *Float_t) parse_frac1(r rune, size int) (next_state next_state_t) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		self.frac_mul *= 10
		self.frac = self.frac*10 + int64(r-'0')
		next_state = self.parse_frac1
	case 'e', 'E':
		next_state = self.parse_exp1
	case 0:
		// ok
	default:
		self.Error = fmt.Sprintf("parse_frac1: invalid format %q", self.input)
	}
	return
}

func (self *Float_t) parse_exp1(r rune, size int) (next_state next_state_t) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		self.Exp = int64(r - '0')
		next_state = self.parse_exp3
	case '-':
		self.exp_sign = true
		next_state = self.parse_exp2
	case '+':
		next_state = self.parse_exp2
	default:
		self.Error = fmt.Sprintf("parse_exp1: invalid format %q", self.input)
	}
	return
}

// here should not be EOF
func (self *Float_t) parse_exp2(r rune, size int) (next_state next_state_t) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		self.Exp *= 10
		self.Exp += int64(r - '0')
		next_state = self.parse_exp3
	default:
		self.Error = fmt.Sprintf("parse_exp2: invalid format %q", self.input)
	}
	return
}

func (self *Float_t) parse_exp3(r rune, size int) (next_state next_state_t) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		self.Exp *= 10
		self.Exp += int64(r - '0')
		if self.Exp > 1024 {
			self.Error = fmt.Sprintf("parse_exp3: overflow %q", self.input)
			return
		}
		next_state = self.parse_exp3
	case 0:
		// ok
	default:
		self.Error = fmt.Sprintf("parse_exp3: invalid format %q", self.input)
	}
	return
}

func (self *Float_t) final() {
	if self.exp_sign {
		self.Exp = -self.Exp
	}
	if self.frac > 0 {
		self.Int = self.Int*self.frac_mul + self.frac
		self.Exp -= CountZero(self.frac_mul)
	}
	if self.int_sign {
		self.Int = -self.Int
	}
}

func ParseFloat(in string) (res Float_t) {
	res.input = in
	res.frac_mul = 1
	next_state := res.parse_int1
	reader := strings.NewReader(in)
	for next_state != nil {
		last_rune, last_size, _ := reader.ReadRune()
		next_state = next_state(last_rune, last_size)
	}
	res.final()
	return
}

func CountZero(in int64) (res int64) {
	for in >= 10 {
		in /= 10
		res++
	}
	return
}

// https://stackoverflow.com/questions/199333/how-do-i-detect-unsigned-integer-overflow

func Add64(a int64, b int64) (int64, bool) {
	if b > 0 && a > math.MaxInt64-b {
		return 0, false
	}
	if b < 0 && a < math.MinInt64-b {
		return 0, false
	}
	return a + b, true
}

func Sub64(a int64, b int64) (int64, bool) {
	if b > 0 && a < math.MinInt64+b {
		return 0, false
	}
	if b < 0 && a > math.MaxInt64+b {
		return 0, false
	}
	return a - b, true
}

func Mul64(a int64, b int64) (int64, bool) {
	if b == 0 {
		return 0, true
	}
	if a > math.MaxInt64/b {
		return 0, false
	}
	if a < math.MinInt64/b {
		return 0, false
	}
	return a * b, true
}
