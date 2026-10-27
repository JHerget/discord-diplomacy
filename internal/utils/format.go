package utils

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ToTitle(value string) string {
	return cases.Title(language.English).String(value)
}
