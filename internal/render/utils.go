package render

import (
	"fmt"
	"strings"
)

// Take a slice of strings and concat each line
func concatLines(strs []string) string {
	var rows []string

	for i := range strs {
		s := strs[i]
		parts := strings.Split(s, "\n")

		for j, part := range parts {
			if len(rows) <= j {
				rows = append(rows, "")
			}

			rows[j] += part
		}
	}

	return strings.Join(rows, "\n")
}

// A really shitty implementation of padLeft
func padLeft(str string, pad int) string {
	return fmt.Sprintf("%"+fmt.Sprint(pad)+"s", str)
}

func reverseString(str string) string {
	lines := strings.Split(str, "\n")
	numLines := len(lines)

	reversed := make([]string, numLines)

	for i, n := range lines {
		j := numLines - i - 1
		reversed[j] = n
	}

	return strings.Join(reversed, "\n")
}
