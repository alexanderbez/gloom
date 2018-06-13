package gloom

import (
	"fmt"
	"hash"
	"hash/fnv"
	"math"

	"github.com/alexanderbez/gloom/murmur3"
)

const (
	// DefaultFalsePosProb is the default value (1%) for the probability of a
	// false positive in a Bloom filter.
	DefaultFalsePosProb = 0.01

	unsetBit bitValue = 0
	setBit   bitValue = 1
)

type (
	// BloomFilter implements a space-efficient randomized data structure for
	// representing a set in order to support membership queries. The filters
	// allow for false positives, however, the space savings often outweigh
	// this drawback. In other words, the filter allows for one-sided errors:
	// Either an element is "probably" in the set or it is definitely not in
	// the set.
	//
	// The Fundamental implementation contains two hash functions that allow
	// for a double hashing technique to facilitate up to k-1 hashes. This is
	// combined with the size of the bit vector, directly correlates to the
	// probability of a false positive.
	//
	// Choice for double hashing was shown to be effective without any loss in
	// the asymptotic false positive probability, leading to less computation and
	// potentially less need for randomness in practice by Adam Kirsch and
	// Michael Mitzenmacher in:
	// 'Less Hashing, Same Performance: Building a Better Bloom Filter'
	//
	// Non-cryptographic hash functions FNV-1a and MurmurHash3 are used for
	// speed performance.
	BloomFilter struct {
		h1, h2    hash.Hash64
		bitVector []bitValue
		m, n, k   uint64
		setBits   uint64
	}

	// bitValue reflects the entry values in the bit vector. It indicates if a
	// set item has potentially been set in the Bloom filter.
	bitValue uint8
)

// NewBloomFilter returns a reference to a new Bloom filter. It requires a,
// potentially approximate, set size n and false positive probability p. The
// optimal value of k, number of hash functions, and m, bit vector size will be
// used.
func NewBloomFilter(n uint64, p float64) (*BloomFilter, error) {
	if n == 0 {
		return nil, fmt.Errorf("invalid set size parameter: %d", n)
	} else if p <= 0 {
		return nil, fmt.Errorf("invalid false positive probability parameter: %f", p)
	}

	m := optimalBitVectorSize(n, p)
	k := optimalNumHash(m, n)

	return &BloomFilter{
		n:         n,
		m:         m,
		k:         k,
		h1:        fnv.New64a(),
		h2:        murmur3.New64(),
		bitVector: make([]bitValue, m, m),
	}, nil
}

// String implements the Stringer interface for a Bloom filter.
func (bf *BloomFilter) String() string {
	return fmt.Sprintf("{n: %d, m: %d, k: %d, setBits: %d}",
		bf.n, bf.m, bf.k, bf.setBits,
	)
}

// NumItemsApprox returns the approximate total number of items in the Bloom
// filter.
func (bf *BloomFilter) NumItemsApprox() uint64 {
	m := float64(bf.m)
	k := float64(bf.k)
	x := float64(bf.setBits)

	return uint64(-(m / k) * math.Log(1-(x/m)))
}

// Includes returns whether or not some arbitrary set item (byte slice) is most
// likely in the Bloom filter. There is a possibility for a false positive with
// the probability being under the Bloom filter's p value. An error is returned
// if any hash function write fails.
func (bf *BloomFilter) Includes(data []byte) (bool, error) {
	if err := bf.hash(data); err != nil {
		return false, err
	}

	var i uint64
	for ; i < bf.k; i++ {
		if bf.bitVector[bf.enhancedDoubleHash(i)] != setBit {
			return false, nil
		}
	}

	return true, nil
}

// Set accepts a set item (byte slice) and sets the appropriate bits to 1 in
// the bit vector. An error is returned if any hash function write fails.
// Enhanced double hashing is used with two hash functions instead of k uniform
// random hash functions.
func (bf *BloomFilter) Set(data []byte) error {
	if err := bf.hash(data); err != nil {
		return err
	}

	var i uint64
	for ; i < bf.k; i++ {
		j := bf.enhancedDoubleHash(i)

		if bf.bitVector[j] == unsetBit {
			bf.setBits++
		}

		bf.bitVector[j] = setBit
	}

	return nil
}

// hash accepts a set item (byte slice) and calculates the two hash values of
// the item. The results are written to the each hash function's internal
// state, so enhanced double hashing can continue in an outside step. Each
// invocation of this call resets each hash function's internal state. An error
// is returned if any hash function write fails.
func (bf *BloomFilter) hash(data []byte) error {
	var err error

	bf.h1.Reset()
	bf.h2.Reset()

	_, err = bf.h1.Write(data)
	if err != nil {
		return err
	}

	_, err = bf.h2.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// enhancedDoubleHash returns the Bloom filter index for a given i such that
// the i is indicative of the kth hash function using enhanced double hashing
// to find the appropriate bit index in the bit vector.
func (bf *BloomFilter) enhancedDoubleHash(i uint64) uint64 {
	g := bf.h1.Sum64() + (i * bf.h2.Sum64()) + uint64(math.Pow(float64(i), 3))
	return g % bf.m
}
