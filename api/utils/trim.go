// Some utility functions
package utils

import (
	"regexp"
)

var trimNewLineRegexp = regexp.MustCompile(" *\r\n *")
var trimAfterBracketsRegexp = regexp.MustCompile("> ")
var trimBeforeBracketsRegexp = regexp.MustCompile(" <")

// TrimNewLines removes leading and trailing line breaks as well as spaces preceding or following line breaks
func TrimNewLines(data []byte) []byte {
	trimmed := trimNewLineRegexp.ReplaceAll(data, []byte(" "))
	trimmed = trimAfterBracketsRegexp.ReplaceAll(trimmed, []byte(">"))
	trimmed = trimBeforeBracketsRegexp.ReplaceAll(trimmed, []byte("<"))
	return trimmed
}
