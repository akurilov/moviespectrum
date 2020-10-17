package util

import (
	"errors"
	"fmt"
	"math"
)

const (
	valueMin          = 0
	valueUpperBound   = 1.0
	medianValue       = valueUpperBound / 2
	medianValueSquare = 0.25
)

func MedianChiSquare(x float64) (float64, error) {
	var err error = nil
	if x < valueMin || x >= valueUpperBound {
		err = errors.New(fmt.Sprintf("The input value should be in the range of [0; 1), got %f", x))
	}
	return (medianValueSquare - math.Pow(medianValue-x, 2)) / medianValueSquare, err
}
