//
//
//

package null

import (
	"io"
	"net/url"
	"strings"
	"unicode/utf8"
)

var StrReplace1 = strings.NewReplacer("'", "''", "\x00", "\\x00", "\x1a", "\\x1a").Replace
var StrReplace2 = strings.NewReplacer("'", "''", "\x00", "\\x00", "\x1a", "\\x1a", "\\", "\\\\").Replace

type StringOption func(in string) string

type LimitWriter_t struct {
	Buf   io.Writer
	Limit int
}

func (self *LimitWriter_t) Write(p []byte) (n int, err error) {
	if self.Limit >= len(p) {
		n, err = self.Buf.Write(p)
	} else {
		for ; self.Limit > 0; self.Limit-- {
			if r, _ := utf8.DecodeLastRune(p[:self.Limit]); r != utf8.RuneError {
				break
			}
		}
		n, err = self.Buf.Write(p[:self.Limit])
	}
	self.Limit -= n
	return
}

func StrLimit(limit int) StringOption {
	var sb strings.Builder
	w := &LimitWriter_t{Buf: &sb, Limit: limit}
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
