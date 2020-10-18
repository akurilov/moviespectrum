package util

import (
	"math"
)

func ChiSquare(x float64, expected float64) float64 {
	expectedSquare := math.Pow(expected, 2)
	return (expectedSquare - math.Pow(expected-x, 2)) / expectedSquare
}
