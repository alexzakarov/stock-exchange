package utils

import (
	"strings"
)

func HtmlMinify(html_data string) string {
	min1 := strings.Replace(html_data, "\t", "", -1)
	min2 := strings.Replace(min1, "\n", "", -1)
	min3 := strings.Replace(min2, " ", "", -1)
	return min3
}
