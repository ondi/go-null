//
//
//

package null

import (
	"fmt"
	"math"
	"strconv"
)

// IEEE 754
// Double: SEEEEEEE EEEEFFFF FFFFFFFF FFFFFFFF FFFFFFFF FFFFFFFF FFFFFFFF FFFFFFFF
// 1 [63]	11 [62–52]	52 [51–00]
// The exponent field is an 11-bit unsigned integer from 0 to 2047, in biased form:
// an exponent value of 1023 represents the actual zero.
// Exponents range from −1022 to +1023 because exponents of −1023 (all 0s) and +1024
// (all 1s) are reserved for special numbers.
// const (
// 	ieee_754_sign = 0b_10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000
// 	ieee_754_expn = 0b_01111111_11110000_00000000_00000000_00000000_00000000_00000000_00000000
// 	ieee_754_frac = 0b_00000000_00001111_11111111_11111111_11111111_11111111_11111111_11111111
// )

func ParseFloatFloat(in float64, frac_overflow bool) (Int int64, Exp int64, err error) {
	if math.IsInf(in, 0) || math.IsNaN(in) {
		err = fmt.Errorf("Inf/Nan not supported")
		return
	}
	return ParseFloatString(strconv.FormatFloat(in, 'f', -1, 64), frac_overflow)
}
