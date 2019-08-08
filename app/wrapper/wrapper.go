package wrapper

import (
	"strings"

	"github.com/eidolon/wordwrap"
)

func Wrap(s string, limit int) string {
	wrapper := wordwrap.Wrapper(limit, false)
	return wrapper(s)
}

func WrapAsComment(s string, limit int) string {
	limit = limit - 3 // we need to account for added comment symbols
	s = "// " + s
	s = Wrap(s, limit)
	s = strings.Replace(s, "\n", "\n// ", -1)

	return s + "\n"
}
