package go2nim

import "strings"

func indent(s string) string {
	return "  " + strings.Replace(s, "\n", "\n  ", -1)
}
