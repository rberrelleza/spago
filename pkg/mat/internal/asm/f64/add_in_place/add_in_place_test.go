package add_in_place

import (
	"fmt"
	// 	"math"
	"math/rand"
	"testing"
)

// const epsilon = 0.00000000000001

func TestAddInPlace(t *testing.T) {
	for length := 0; length < 1000; length++ {
		t.Run(fmt.Sprintf("length %d", length), func(t *testing.T) {
			x := randomVector(length)
			xCopy := copySlice(x)
			yActual := randomVector(length)
			yExpected := copySlice(yActual)

			AddInPlace(x, yActual)
			if !sliceEqual(x, xCopy) {
				t.Error("x was modified")
			}
			goAddInPlace(x, yExpected)

			if !sliceEqual(yActual, yExpected) {
				t.Errorf("expected %v, actual %v", yExpected, yActual)
			}
		})
	}
}

func goAddInPlace(x, y []float64) {
	for i, xVal := range x {
		y[i] += xVal
	}
}

func randomVector(length int) []float64 {
	x := make([]float64, length)
	for i := range x {
		x[i] = rand.Float64() * 100
	}
	return x
}

func copySlice(x []float64) []float64 {
	y := make([]float64, len(x))
	copy(y, x)
	return y
}

func sliceEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i, aVal := range a {
		if !floatEqual(aVal, b[i]) {
			return false
		}
	}
	return true
}

func floatEqual(a, b float64) bool {
	const delta = 0.0000000000001
	return b > a-delta && b < a+delta

}
