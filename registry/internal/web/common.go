package web

import (
	"bufio"
	"strings"
)

func linesOfText(text string) []string {
	lines := make([]string, 0, 1)
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, scanner.Text())
		}
	}
	return lines
}
