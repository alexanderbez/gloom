package gloom

import "math"

var (
	lnTwoSq = math.Pow(math.Log(2), 2)
	lnTwo   = math.Log(2)
)

// optimalBitVectorSize returns the optimal bit vector size for a Bloom filter
// given a set size n and a false positive probability p.
func optimalBitVectorSize(n uint, p float64) uint {
	return uint(math.Ceil(-((float64(n) * math.Log(p)) / lnTwoSq)))
}

// optimalNumHash returns the optimal number of hash 'functions' for a Bloom
// filter given a bit vector size m and a set size n.
func optimalNumHash(m, n uint) uint {
	return uint(math.Ceil(float64(m/n) * lnTwo))
}
