package sanitize

import (
	"strings"
)

func Name(name string) string {
	splittedName := strings.FieldsFunc(name, split)
	var sanitizedName string
	for _, splitted := range splittedName {
		sanitizedName = sanitizedName + strings.Title(splitted)
	}

	return strings.TrimSpace(sanitizedName)
}

func split(r rune) bool {
	return r == '-' || r == '_'
}
