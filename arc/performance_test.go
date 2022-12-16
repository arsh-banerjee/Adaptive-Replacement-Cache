/******************************************************************************
Filename: performance_test.go
Names: Arsh Banerjee, Kenny Lam, and Kenar Vyas
NetId:
Description:
*****************************************************************************/

package arc

import (
	"math/rand"
	"strconv"
	"testing"
)

func TestCacheBenchmark(t *testing.T) {
	capacity := 1000
	cache := NewArc(capacity)

	rand.Seed(2)

	tracesize := capacity * 10
	var trace [10000]int
	for i := 0; i < tracesize; i++ {
		trace[i] = rand.Intn(tracesize / 2)
		cache.Set(strconv.Itoa(trace[i]), []byte(""))
	}

	hits := 0
	misses := 0

	for i := 0; i < tracesize; i++ {
		trace[i] = rand.Intn(tracesize / 2)
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
