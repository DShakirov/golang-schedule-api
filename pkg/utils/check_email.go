package utils

import (
	"regexp"
)

func IsValidEmail(email string) bool {
	// Using regex for validating email
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}
