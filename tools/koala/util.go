package main

import (
	"strings"
	"text/template"
)

var templateFuncMap = template.FuncMap{
	"Capitalize": Capitalize,
}

func Capitalize(str string) string {
	var output string
	chars := []rune(str)
	for index, char := range chars {
		if index == 0 {
			if char < 'a' || char > 'z' {
				return output
			}
			output += strings.ToUpper(string(char))
			continue
		}
		output += string(char)
	}

	return output
}
