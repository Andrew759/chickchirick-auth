package service

import (
	"regexp"
	"strconv"
	"unicode/utf8"
)

func IsLatinSymbolOnly(s string) bool {
	return regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(s)
}

func IsLogin(s string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_]{2,256}$`).MatchString(s)
}

func IsPhoneNumber(s string) bool {
	return regexp.MustCompile(`^\+?[1-9]\d{6,14}$`).MatchString(s)
}

func IsEmail(s string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`).MatchString(s)
}

func IsHasCorrectLength(s string, maxLength int) bool {
	sLen := utf8.RuneCountInString(s)

	return sLen > 0 && sLen <= maxLength
}

func isInteger(s string) bool {
	_, err := strconv.Atoi(s)

	return err == nil
}
