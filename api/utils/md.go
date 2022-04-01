package utils

import (
	"bytes"

	"github.com/yuin/goldmark"
)

func Md2Html(md string) (string, error) {
	var html bytes.Buffer
	if err := goldmark.Convert([]byte(md), &html); err != nil {
		return "", err
	}
	return string(html.Bytes()), nil
}
