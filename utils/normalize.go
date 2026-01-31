package utils

import (
	"strconv"
	"strings"
	"time"
)

// NormalizeYear converts string or numeric year to int
func NormalizeYear(year interface{}) int {
	switch v := year.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		yearStr := strings.TrimSpace(v)
		y, err := strconv.Atoi(yearStr)
		if err != nil {
			return 0
		}
		return y
	default:
		return 0
	}
}

// normalizeDate converts various formats → YYYY-MM-DD
func NormalizeDate(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}

	layouts := []string{
		"2006-01-02",
		"02 Jan 2006",
		"January 2, 2006",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, v); err == nil {
			return t.Format("2006-01-02")
		}
	}

	return v // fallback, better than losing data
}

// NormalizeLanguage trims and lowercases language string
func NormalizeLanguage(v any) string {
	switch t := v.(type) {
	case string:
		return strings.TrimSpace(t)
	case []any:
		langs := []string{}
		for _, l := range t {
			if s, ok := l.(string); ok {
				langs = append(langs, strings.TrimSpace(s))
			}
		}
		return strings.Join(langs, ", ")
	default:
		return ""
	}
}

// NormalizeSlice trims each element of a string slice
func NormalizeSlice(in []string) []string {
	out := []string{}
	for _, v := range in {
		v = strings.TrimSpace(v)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}
