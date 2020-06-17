package wrrimpl

import (
	"math"
	"testing"

	"github.com/c2fo/testify/assert"
)

func assertApproxWithPresicion(t *testing.T, iter int, got, want, presicion float64) {
	t.Helper()

	delta := math.Abs(got - want)
	mean := math.Abs(got+want) / 2.
	assert.True(t, delta < 1e-6 || delta/mean < presicion, "iter: %d, want: %0.4f, got: %0.4f, precision: %0.4f", iter, want, got, presicion)
}

func testSuiteNext(t *testing.T, f func() WRR, p float64) {
	const invocationsCount = 1000

	// NOTE: there is an edge case for 1%, 99% distribution
	// in this case actual diff in p.p. is 0.1 (want 1%, got 1.1%)
	// but for this case the basic deviation formula should be changed
	// i don't mind this case tbh (especially for that small amount of queries)

	for _, tc := range []struct {
		name    string
		weights []int64
	}{
		{
			name:    "1-1-1",
			weights: []int64{1, 1, 1},
		},
		{
			name:    "1-2-3",
			weights: []int64{1, 2, 3},
		},
		{
			name:    "5-3-2",
			weights: []int64{5, 3, 2},
		},
		{
			name:    "17-23-37",
			weights: []int64{17, 23, 37},
		},
		{
			name:    "80-20",
			weights: []int64{80, 20},
		},
		{
			name:    "100-0",
			weights: []int64{100, 0},
		},
		{
			name:    "20-40-20-40",
			weights: []int64{20, 40, 20, 40},
		},
		{
			name:    "50-50",
			weights: []int64{50, 50},
		},
		{
			name:    "10-90",
			weights: []int64{10, 90},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var sumOfWeights int64

			w := f()
			for i, weight := range tc.weights {
				w.Add(i, weight)
				sumOfWeights += weight
			}

			results := make(map[int]int)
			for i := 0; i < invocationsCount; i++ {
				results[w.Next().(int)]++
			}

			wantRatio := make([]float64, len(tc.weights))
			for i, weight := range tc.weights {
				wantRatio[i] = float64(weight) / float64(sumOfWeights)
			}
			gotRatio := make([]float64, len(tc.weights))
			for i, c := range results {
				gotRatio[i] = float64(c) / invocationsCount
			}

			for i := range wantRatio {
				assertApproxWithPresicion(t, i, gotRatio[i], wantRatio[i], p)
			}
		})
	}
}

func TestEDF(t *testing.T) {
	testSuiteNext(t, NewEDF, 0.01)
}
func TestRandom(t *testing.T) {
	// random distribution is much worse and not so predictable
	// even with EDF's deadline possible interchange problems
	testSuiteNext(t, NewRandom, 0.5)
}
