//
//
//

package null

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"
)

var StrReplace1 = strings.NewReplacer("'", "''", "\x00", "\\x00", "\x1a", "\\x1a").Replace
var StrReplace2 = strings.NewReplacer("'", "''", "\x00", "\\x00", "\x1a", "\\x1a", "\\", "\\\\").Replace

type StringOption func(in string) string

type Limit_t struct {
	Bytes int
}

func (self Limit_t) Limit(in string) string {
	if len(in) > self.Bytes {
		for ; self.Bytes > 0; self.Bytes-- {
			if r, _ := utf8.DecodeLastRuneInString(in[:self.Bytes]); r != utf8.RuneError {
				break
			}
		}
		return in[:self.Bytes]
	}
	return in
}

func StrLimit(limit int) StringOption {
	return Limit_t{Bytes: limit}.Limit
}

func StrQuote1(in string) string {
	return "'" + StrReplace1(in) + "'"
}

func StrQuote2(in string) string {
	return "'" + StrReplace2(in) + "'"
}

func PowInt64(x int64, n int64) (res int64) {
	if n == 0 {
		return 1
	}
	if n == 1 {
		return x
	}
	res = PowInt64(x, n/2)
	if n%2 == 0 {
		return res * res
	}
	return x * res * res
}

func Degree(in int64, by int64) (res int64) {
	for in /= by; in != 0; in /= by {
		res++
	}
	return
}

func LeadZero(in string) (res int64) {
	for _, v := range in {
		if v != '0' {
			return
		}
		res++
	}
	return
}

func StringPriceToInt64(in string, mul int64) (res int64, err error) {
	ix := strings.Index(in, ".")
	if ix == -1 {
		res, err = strconv.ParseInt(in, 10, 64)
		res = res * mul
		return
	}
	if res, err = strconv.ParseInt(in[:ix], 10, 64); err != nil {
		return
	}
	frac, err := strconv.ParseInt(in[ix+1:], 10, 64)
	if err != nil {
		return
	}
	if frac < 0 {
		err = errors.New("fraction format error")
		return
	}
	shift := Degree(mul, 10) - Degree(frac, 10) - LeadZero(in[ix+1:]) - 1
	if shift < 0 {
		frac /= PowInt64(10, -shift)
	} else {
		frac *= PowInt64(10, shift)
	}
	if res < 0 {
		res = res*mul - frac
	} else {
		res = res*mul + frac
	}
	return
}
