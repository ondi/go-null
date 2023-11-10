//
//
//

package null

import (
	"errors"
	"io"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"
)

var StrReplace1 = strings.NewReplacer("'", "''", "\x00", "\\x00", "\x1a", "\\x1a").Replace
var StrReplace2 = strings.NewReplacer("'", "''", "\x00", "\\x00", "\x1a", "\\x1a", "\\", "\\\\").Replace

type StringOption func(in string) string

type LimitWriter_t struct {
	Out     io.Writer
	Limit   int
	written int
}

func (self *LimitWriter_t) Write(p []byte) (n int, err error) {
	if n = self.Limit - self.written; n > len(p) {
		n, err = self.Out.Write(p)
	} else {
		for ; n > 0; n-- {
			if r, _ := utf8.DecodeLastRune(p[:n]); r != utf8.RuneError {
				break
			}
		}
		n, err = self.Out.Write(p[:n])
	}
	self.written += n
	return
}

func StrLimit(limit int) StringOption {
	var sb strings.Builder
	w := &LimitWriter_t{Out: &sb, Limit: limit}
	return func(in string) string {
		io.WriteString(w, in)
		return sb.String()
	}
}

func StrQuote1(in string) string {
	return "'" + StrReplace1(in) + "'"
}

func StrQuote2(in string) string {
	return "'" + StrReplace2(in) + "'"
}

func StrUrlEscape(in string) string {
	return url.QueryEscape(in)
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
