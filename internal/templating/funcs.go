package templating

import (
	"strings"
	"text/template"
	"unicode"
)

// FuncMap returns the custom template functions available to every template
// rendered by the Engine.
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"ToCamel":      ToCamel,
		"ToSnake":      ToSnake,
		"ToLowerFirst": ToLowerFirst,
		"ToPlural":     ToPlural,
	}
}

// ToCamel converts "widget_order" or "widget-order" or "widget order" to "WidgetOrder".
func ToCamel(s string) string {
	fields := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})
	var b strings.Builder
	for _, f := range fields {
		if f == "" {
			continue
		}
		r := []rune(f)
		b.WriteRune(unicode.ToUpper(r[0]))
		b.WriteString(strings.ToLower(string(r[1:])))
	}
	return b.String()
}

// ToSnake converts "WidgetOrder" or "widget order" to "widget_order".
func ToSnake(s string) string {
	var b strings.Builder
	runes := []rune(s)
	for i, r := range runes {
		if r == '-' || r == ' ' {
			b.WriteRune('_')
			continue
		}
		if unicode.IsUpper(r) {
			if i > 0 && (unicode.IsLower(runes[i-1]) || unicode.IsDigit(runes[i-1])) {
				b.WriteRune('_')
			}
			b.WriteRune(unicode.ToLower(r))
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// ToLowerFirst lower-cases the first rune of s, leaving the rest untouched.
func ToLowerFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

// ToPlural is a simple suffix-based English pluralizer. It covers the common
// cases well enough for generated resource names; it is not a full
// inflection engine by design (kept intentionally simple for MVP).
func ToPlural(s string) string {
	if s == "" {
		return s
	}
	lower := strings.ToLower(s)
	switch {
	case strings.HasSuffix(lower, "y") && len(s) > 1 && !isVowel(rune(lower[len(lower)-2])):
		return s[:len(s)-1] + "ies"
	case strings.HasSuffix(lower, "s"),
		strings.HasSuffix(lower, "x"),
		strings.HasSuffix(lower, "z"),
		strings.HasSuffix(lower, "ch"),
		strings.HasSuffix(lower, "sh"):
		return s + "es"
	default:
		return s + "s"
	}
}

func isVowel(r rune) bool {
	switch unicode.ToLower(r) {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	}
	return false
}
