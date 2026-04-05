package internal

import (
	"strings"
	"text/template"
	"unicode"
)

// StringFuncs returns a map of string manipulation functions for templates.
func StringFuncs() template.FuncMap {
	return template.FuncMap{
		"toLower":          strings.ToLower,
		"toUpper":          strings.ToUpper,
		"capitalize":       capitalize,
		"toCamelCase":      toCamelCase,
		"toPascalCase":     toPascalCase,
		"toSnakeCase":      toSnakeCase,
		"toKebabCase":      toKebabCase,
		"toScreamingSnake": toScreamingSnake,
		"toScreamingKebab": toScreamingKebab,
		"toDotCase":        toDotCase,
		"toPathCase":       toPathCase,
		"trimSpace":        strings.TrimSpace,
		"hasPrefix":        strings.HasPrefix,
		"hasSuffix":        strings.HasSuffix,
		"contains":         strings.Contains,
		"replace":          strings.ReplaceAll,
		"split":            strings.Split,
		"join":             strings.Join,
	}
}

// capitalize capitalizes the first character of the string and lowers the rest.
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	r := []rune(s)
	return string(unicode.ToUpper(r[0])) + strings.ToLower(string(r[1:]))
}

// splitWords splits a string into words based on common delimiters and case transitions.
// It handles snake_case, kebab-case, camelCase, PascalCase, dot.case, path/case, and spaces.
func splitWords(s string) []string {
	var words []string
	var current []rune
	runes := []rune(s)

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		// Delimiters
		if r == '_' || r == '-' || r == '.' || r == '/' || unicode.IsSpace(r) {
			if len(current) > 0 {
				words = append(words, string(current))
				current = nil
			}
			continue
		}

		// Word boundary: lower/digit to upper
		if i > 0 && unicode.IsUpper(r) && (unicode.IsLower(runes[i-1]) || unicode.IsDigit(runes[i-1])) {
			if len(current) > 0 {
				words = append(words, string(current))
				current = nil
			}
		}

		// Word boundary: upper to upper+lower (e.g., XMLParser -> XML Parser)
		if i > 0 && unicode.IsUpper(r) && unicode.IsUpper(runes[i-1]) {
			if i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
				if len(current) > 0 {
					words = append(words, string(current))
					current = nil
				}
			}
		}

		current = append(current, r)
	}

	if len(current) > 0 {
		words = append(words, string(current))
	}

	return words
}

func toCamelCase(s string) string {
	words := splitWords(s)
	for i, word := range words {
		if i == 0 {
			words[i] = strings.ToLower(word)
		} else {
			words[i] = capitalize(word)
		}
	}
	return strings.Join(words, "")
}

func toPascalCase(s string) string {
	words := splitWords(s)
	for i, word := range words {
		words[i] = capitalize(word)
	}
	return strings.Join(words, "")
}

func toSnakeCase(s string) string {
	words := splitWords(s)
	for i, word := range words {
		words[i] = strings.ToLower(word)
	}
	return strings.Join(words, "_")
}

func toKebabCase(s string) string {
	words := splitWords(s)
	for i, word := range words {
		words[i] = strings.ToLower(word)
	}
	return strings.Join(words, "-")
}

func toScreamingSnake(s string) string {
	words := splitWords(s)
	for i, word := range words {
		words[i] = strings.ToUpper(word)
	}
	return strings.Join(words, "_")
}

func toScreamingKebab(s string) string {
	words := splitWords(s)
	for i, word := range words {
		words[i] = strings.ToUpper(word)
	}
	return strings.Join(words, "-")
}

func toDotCase(s string) string {
	words := splitWords(s)
	for i, word := range words {
		words[i] = strings.ToLower(word)
	}
	return strings.Join(words, ".")
}

func toPathCase(s string) string {
	words := splitWords(s)
	for i, word := range words {
		words[i] = strings.ToLower(word)
	}
	return strings.Join(words, "/")
}
