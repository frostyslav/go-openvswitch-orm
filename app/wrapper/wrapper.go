package wrapper

import (
	"fmt"
	"strings"

	"github.com/eidolon/wordwrap"
)

func Wrap(s string) string {
	wrapper := wordwrap.Wrapper(77, false)
	return wrapper(s)
}

func WrapAsComment(s string) string {
	s = fmt.Sprintf("// %s", s)
	s = Wrap(s)
	s = strings.Replace(s, "\n", "\n// ", -1)

	return fmt.Sprintf("%s\n", s)
}
