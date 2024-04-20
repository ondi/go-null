//
//
//

package null

import (
	"fmt"
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
	Input     string
	IntPart   int64
	frac_part int64
	frac_exp  int64
	Exp       int64
	int_sign  bool
	exp_sign  bool
	Error     string
}

func (self *Float_t) state_int_part_first(r rune, size int) (next_state next_state_t) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		self.IntPart = int64(r - '0')
		next_state = self.state_int_part_third
	case '-':
		self.int_sign = true
		next_state = self.state_int_part_second
	case '+':
		next_state = self.state_int_part_second
	case '.':
		next_state = self.state_frac_part_first
	default:
		self.Error = fmt.Sprintf("state_int_part_first: unexpected %q", r)
	}
	return
}

// here should not be EOF
func (self *Float_t) state_int_part_second(r rune, size int) (next_state next_state_t) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		self.IntPart = self.IntPart*10 + int64(r-'0')
		next_state = self.state_int_part_third
	case '.':
		next_state = self.state_frac_part_first
	default:
		self.Error = fmt.Sprintf("state_int_part_second: unexpected %q", r)
	}
	return
}

func (self *Float_t) state_int_part_third(r rune, size int) (next_state next_state_t) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		self.IntPart = self.IntPart*10 + int64(r-'0')
		next_state = self.state_int_part_third
	case '.':
		next_state = self.state_frac_part_first
	case 0:
		// ok
	default:
		self.Error = fmt.Sprintf("state_int_part_third: unexpected %q", r)
	}
	return
}

func (self *Float_t) state_frac_part_first(r rune, size int) (next_state next_state_t) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		self.frac_exp++
		self.frac_part = self.frac_part*10 + int64(r-'0')
		next_state = self.state_frac_part_first
	case 'e', 'E':
		next_state = self.state_exponent_first
	case 0:
		// ok
	default:
		self.Error = fmt.Sprintf("state_frac_part_first: unexpected %q", r)
	}
	return
}

func (self *Float_t) state_exponent_first(r rune, size int) (next_state next_state_t) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		self.Exp = int64(r - '0')
		next_state = self.state_exponent_third
	case '-':
		self.exp_sign = true
		next_state = self.state_exponent_second
	case '+':
		next_state = self.state_exponent_second
	default:
		self.Error = fmt.Sprintf("state_exponent_first: unexpected %q", r)
	}
	return
}

// here should not be EOF
func (self *Float_t) state_exponent_second(r rune, size int) (next_state next_state_t) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		self.Exp *= 10
		self.Exp += int64(r - '0')
		next_state = self.state_exponent_third
	default:
		self.Error = fmt.Sprintf("state_exponent_second: unexpected %q", r)
	}
	return
}

func (self *Float_t) state_exponent_third(r rune, size int) (next_state next_state_t) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		self.Exp *= 10
		self.Exp += int64(r - '0')
		next_state = self.state_exponent_third
	case 0:
		// ok
	default:
		self.Error = fmt.Sprintf("state_exponent_third: unexpected %q", r)
	}
	return
}

func (self *Float_t) final() {
	if self.exp_sign {
		self.Exp = -self.Exp
	}
	if self.frac_part > 0 {
		self.IntPart = self.IntPart*Width10(self.frac_exp) + self.frac_part
		self.Exp -= self.frac_exp
	}
	if self.int_sign {
		self.IntPart = -self.IntPart
	}
}

func ParseFloat(in string) (res Float_t) {
	res.Input = in
	next_state := res.state_int_part_first
	reader := strings.NewReader(in)
	for next_state != nil {
		last_rune, last_size, _ := reader.ReadRune()
		next_state = next_state(last_rune, last_size)
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
