package datafunctions

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
)

// This function will calculate the correlation between two variables
// Pearson's correlation coefficient
// r = n(∑xy) - (∑x)(∑y) / √[n(∑x^2) - (∑x)^2][n(∑y^2) - (∑y)^2]
func Correlation(csvFileName string) float64 {
	file, err := os.Open(csvFileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	numberOfDataPoints := len(data) - 1 //n value
	//x and y values
	sumOfXValues := 0.0
	sumOfX2Values := 0.0
	sumOfYValues := 0.0
	sumOfY2Values := 0.0
	sumOfXYValues := 0.0
	for v := 1; v < len(data); v++ {
		x, err := strconv.ParseFloat(data[v][1], 64)
		if err != nil {
			fmt.Println(err)
		}
		sumOfXValues += x
		sumOfX2Values += (x * x)
		y, _ := strconv.ParseFloat(data[v][2], 64)
		sumOfYValues += y
		sumOfY2Values += (y * y)
		sumOfXYValues += x * y
	}
	// Calculate r
	a := (float64(numberOfDataPoints) * sumOfXYValues) - (sumOfXValues * sumOfYValues)
	b := (float64(numberOfDataPoints) * sumOfX2Values) - (sumOfXValues * sumOfXValues)
	c := (float64(numberOfDataPoints) * sumOfY2Values) - (sumOfYValues * sumOfYValues)

	r := a / math.Sqrt(b*c)

	return r
}
