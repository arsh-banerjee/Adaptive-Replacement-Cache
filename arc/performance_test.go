/******************************************************************************
Filename: performance_test.go
Names: Arsh Banerjee, Kenny Lam, and Kenar Vyas
NetId: arshb, kennyl, kvyas
Description: Performance test to compare ARC and LRU Caches
*****************************************************************************/

package arc

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"testing"
)

// Helper function copied from: https://stackoverflow.com/questions/18390266/how-can-we-truncate-float64-type-to-a-particular-precision
func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

// Normal Distribution (Not realistic workload)
func TestCacheBenchmark(t *testing.T) {
	capacity := 10000
	cache := NewArc(capacity)

	rand.Seed(2)

	tracesize := 100000
	var trace [100000]int
	for i := 0; i < tracesize; i++ {
		trace[i] = rand.Intn(5000)
		cache.Set(strconv.Itoa(trace[i]), []byte(""))
	}

	hits := 0
	misses := 0

	for i := 0; i < tracesize; i++ {
		trace[i] = rand.Intn(5000)
		_, prs := cache.Get(strconv.Itoa(trace[i]))
		if prs {
			hits += 1
		} else {
			misses += 1
		}
	}

	ratio := float64(hits) / (float64(hits) + float64(misses))
	t.Logf("ARC of size %d: hits: %d | misses: %d | hit-ratio: %f", capacity, hits, misses, ratio)

}

//Exponential Distribution (More realistic web workload - some are visited way more than others)
func TestCacheBenchmark2(t *testing.T) {
	capacities := []int{10000, 12000, 14000, 16000, 18000, 20000}

	for i := 0; i < len(capacities); i++ {
		capacity := capacities[i]
		cache := NewArc(capacity)

		rand.Seed(2)

		tracesize := 100000
		var trace [100000]float64
		for i := 0; i < tracesize; i++ {
			trace[i] = toFixed(rand.ExpFloat64(), 3)
			s := fmt.Sprintf("%f", trace[i])
			cache.Set(s, []byte(""))
		}

		hits := 0
		misses := 0

		for i := 0; i < tracesize; i++ {
			trace[i] = toFixed(rand.ExpFloat64(), 3)
			s := fmt.Sprintf("%f", trace[i])
			_, prs := cache.Get(s)
			if prs {
				hits += 1
			} else {
				misses += 1
			}
		}

		ratio := float64(hits) / (float64(hits) + float64(misses))
		t.Logf("ARC of size %d: hits: %d | misses: %d | hit-ratio: %f", capacity, hits, misses, ratio)
	}
}
