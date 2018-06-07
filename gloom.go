package gloom

import (
	"hash"
	"hash/fnv"
)

// REFS:
//
// https://en.wikipedia.org/wiki/Bloom_filter
// http://citeseerx.ist.psu.edu/viewdoc/download;jsessionid=AF0A7F109B5F97C758DD773942A1054F?doi=10.1.1.152.579&rep=rep1&type=pdf
// https://github.com/jedisct1/rust-bloom-filter/blob/master/src/bloomfilter/lib.rs
// https://llimllib.github.io/bloomfilter-tutorial/

// h := fnv.New32a()

// io.WriteString(h, "The fog is getting thicker!")
// fmt.Println(h.Sum32() % 1024)
// h.Reset()

const (
	// DefaultFalseNegProb is the default value for the probability of a false
	// negative in a Bloom filter.
	DefaultFalseNegProb = 0.01
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
		fnvHasher hash.Hash32
		bitVector []byte
		m, n, k   uint
	}
)

// NewBloomFilter returns a reference to a new Bloom filter. It requires a,
// potentially approximate, set size n and false positive probability p. The
// optimal value of k, number of hash functions, and m, bit vector size will be
// used.
func NewBloomFilter(n uint, p float64) *BloomFilter {
	m := optimalBitVectorSize(n, p)
	k := optimalNumHash(m, n)

	return &BloomFilter{
		n:         n,
		m:         m,
		k:         k,
		fnvHasher: fnv.New32a(),
		bitVector: make([]byte, m, m),
	}
}

// fnvHash returns a 32 bit FNV-1a hash of a given slice. An error is returned
// if the byte slice cannot be written to the underlying FNV-1a hash writer.
func (bf *BloomFilter) fnvHash(b []byte) (uint32, error) {
	defer bf.fnvHasher.Reset()

	if _, err := bf.fnvHasher.Write(b); err != nil {
		return 0, err
	}

	return bf.fnvHasher.Sum32(), nil
}
