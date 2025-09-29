package parser

import "strings"

func Parse(argStr string) []string {
	return strings.Fields(argStr)
}
