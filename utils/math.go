package utils

import (
	"math"
	"regexp"
	"strconv"
)

var (
	// MaxInt max int value
	MaxInt int
	// MaxUint max uint value
	MaxUint uint
	// MinInt min int value
	MinInt int
)

func init() {
	switch strconv.IntSize {
	case 64:
		MaxInt = math.MaxInt64
		MinInt = math.MinInt64
		MaxUint = math.MaxUint64
	case 32:
		MaxInt = math.MaxInt32
		MinInt = math.MinInt32
		MaxUint = math.MaxUint32
	case 16:
		MaxInt = math.MaxInt16
		MinInt = math.MinInt16
		MaxUint = math.MaxUint16
	case 8:
		MaxInt = math.MaxInt8
		MinInt = math.MinInt8
		MaxUint = math.MaxUint8
	}
}

var matchName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]+$`)
var matchPassword = regexp.MustCompile(`^[a-f0-9]+$`)

// MatchName if match user name return true
func MatchName(val string) bool {
	if len(val) < 4 {
		return false
	}
	return matchName.MatchString(val)
}

// MatchPassword if match passsword name return true
func MatchPassword(val string) bool {
	if len(val) != 32 {
		return false
	}
	return matchPassword.MatchString(val)
}
