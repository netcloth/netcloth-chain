package utils

import "regexp"

var IsAlphaNumeric = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
