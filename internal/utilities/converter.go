package utilities

import (
	"math"
	"strings"
)

func NewBoolean(value bool) *bool {
	b := value
	return &b
}

func NewString(value string) *string {
	s := value
	return &s
}

func NewFloat(value float64) *float64 {
	s := value
	return &s
}

func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(Round(num*output)) / output
}

func AddingZeroInFront(value string) string {
	length := len(value)
	if length < 2 {
		value = "0" + value
	}
	return value
}

func RemoveZeroInFront(value string) string {
	length := len(value)
	if length >= 2 {
		value = strings.ReplaceAll(value, "0", "")
	}
	return value
}
