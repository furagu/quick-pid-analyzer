package utils

import (
	"fmt"
	"strconv"
	"strings"
)

const floatsLimit = 6
const floatsDelimiter = ", "

func FloatsToString(f []float64) string {
	if len(f) < floatsLimit {
		return strings.Join(floatsToStringList(f), floatsDelimiter)
	} else {
		return fmt.Sprintf(
			"[%s ... %s]",
			strings.Join(floatsToStringList(f[:floatsLimit/2]), floatsDelimiter),
			strings.Join(floatsToStringList(f[len(f)-floatsLimit/2:]), floatsDelimiter),
		)
	}
}

func floatsToStringList(f []float64) []string {
	s := make([]string, len(f), len(f))
	for i, v := range f {
		s[i] = strconv.FormatFloat(v, 'f', -1, 64)
	}
	return s
}
