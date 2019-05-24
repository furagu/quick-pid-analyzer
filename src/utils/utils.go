package utils

import (
	"fmt"
	"strconv"
	"strings"
)

const floatsLimit = 6
const floatsDelimiter = ", "

func Ppaa(f [][]float64) string {
	s := []string{}
	if len(f) < floatsLimit {
		for i, v := range f {
			s = append(s, fmt.Sprintf("%d: %s", i, FloatsToString(v)))
		}
	} else {
		for i, v := range f[0 : floatsLimit/2] {
			s = append(s, fmt.Sprintf("%d: %s", i, FloatsToString(v)))
		}
		s = append(s, "...")
		for i, v := range f[len(f)-floatsLimit/2 : len(f)] {
			s = append(s, fmt.Sprintf("%d: %s", len(f)-floatsLimit/2+i, FloatsToString(v)))
		}
	}
	return strings.Join(s, "\n")
}

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
