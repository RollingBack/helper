package helper

import "regexp"

func IsHexString(testString string) bool {
	reg := regexp.MustCompile(`(\\x[0-9a-f]{1,2}])*`)
	return reg.MatchString(testString)
}
