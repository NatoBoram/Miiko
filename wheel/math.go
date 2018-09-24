package wheel

import (
	"math"
	"math/rand"
	"time"
)

// Phi represents the Golden Ratio.
func Phi() float64 {
	return (1 + math.Sqrt(5)) / 2
}

// MaxInt returns the larger of x or y.
func MaxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// MinInt returns the smaller of x or y.
func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// StringInSlice checks if a string is inside a slice.
func StringInSlice(str string, list []string) bool {
	for _, item := range list {
		if item == str {
			return true
		}
	}
	return false
}

// RandomOverPhiPower gives a random `bool` according to the following formula : `random <= 1/phi^x`.
// The larger is `x`, the smaller the chances of obtaining `true`.
//
// Examples :
// 0 = 100%
// 1 =  62%
// 2 =  38%
// 3 =  24%
// 5 =   9%
// 10 =  1%
// 100 = 0%
func RandomOverPhiPower(power float64) bool {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Float64() <= 1/math.Pow(Phi(), power)
}
