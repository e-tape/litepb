package generator

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var casesTitleRegex = regexp.MustCompile("([A-Z])")

func generateName(value string) string {
	pName := strings.TrimLeft(value, "_")
	prefix := strings.TrimSuffix(value, pName)
	pName = strings.ReplaceAll(
		cases.Title(language.English).String(
			strings.ReplaceAll(
				casesTitleRegex.ReplaceAllString(pName, " $1"),
				"_",
				" ",
			),
		),
		" ",
		"",
	)
	pName = prefix + pName
	pName, ok := strings.CutPrefix(pName, "__")
	if ok {
		pName = "X" + pName
	}
	pName, ok = strings.CutPrefix(pName, "_")
	if ok {
		pName = "X" + pName
	}
	return pName
}
