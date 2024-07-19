package helpers

import "math"

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(Round(num*output)) / output
}

func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func Average(numbers []float64) float64 {
	size := len(numbers)
	// add all numbers in array
	total := 0.0
	for _, number := range numbers {
		total += number
	}

	return total / float64(size)
}
