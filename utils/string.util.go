package utils

import (
	"path/filepath"
	"strings"
)

var uppercaseAcronym = map[string]string{
	"ID": "id",
}

// ConfigureAcronym allows you to add additional words which will be considered acronyms
func ConfigureAcronym(key, val string) {
	uppercaseAcronym[key] = val
}
func toCamelInitCase(s string, initCase bool) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	if a, ok := uppercaseAcronym[s]; ok {
		s = a
	}

	n := strings.Builder{}
	n.Grow(len(s))
	capNext := initCase
	for i, v := range []byte(s) {
		vIsCap := v >= 'A' && v <= 'Z'
		vIsLow := v >= 'a' && v <= 'z'
		if capNext {
			if vIsLow {
				v += 'A'
				v -= 'a'
			}
		} else if i == 0 {
			if vIsCap {
				v += 'a'
				v -= 'A'
			}
		}
		if vIsCap || vIsLow {
			n.WriteByte(v)
			capNext = false
		} else if vIsNum := v >= '0' && v <= '9'; vIsNum {
			n.WriteByte(v)
			capNext = true
		} else {
			capNext = v == '_' || v == ' ' || v == '-' || v == '.'
		}
	}
	return n.String()
}

// ToCamel converts a string to CamelCase
func ToCamel(s string) string {
	return toCamelInitCase(s, true)
}

// ToLowerCamel converts a string to lowerCamelCase
func ToLowerCamel(s string) string {
	return toCamelInitCase(s, false)
}

// need to write test for this function, not work properly
func HasSuffixCheck(name string, checks []string) bool {
	for _, check := range checks {
		if !strings.HasSuffix(name, "."+check) {
			return false
		}
	}
	return true
}
func HasImageSuffixCheck(name string) bool {
	imageExtensions := []string{"jpg", "jpeg", "png", "gif", "bmp", "tif", "tiff", "webp"}
	return HasSuffixCheck(name, imageExtensions)
}
func IsImageFile(fileName string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tif", ".tiff", ".webp":
		return true
	default:
		return false
	}
}
func SplitString(input string, separator ...string) []string {
	if input == "" {
		return nil
	}
	var stringSeperator string
	if separator[0] == "" {
		stringSeperator = ","
	}

	return strings.Split(input, stringSeperator)
}

// returns xth of english alphabet letter in upper case
func GetNthAlphabeticUpperLetter(x int) string {
	if x < 1 || x > 26 {
		return "" // Invalid input range, return an empty string or handle error as needed
	}

	letter := 'A' + rune(x-1)
	return string(letter)
}
