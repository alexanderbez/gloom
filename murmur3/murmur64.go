package murmur3

import "hash"

type (
	// sum64 implements the non-cryptographic Murmur3 64-bit hash function
	// suitable for general hash-based lookup.
	//
	// NOTE: There is no direct Murmur3 64-bit variant protocol, so the 128-bit
	// variant is used where only the most significant 64-bit unsigned integer
	// is used.
	sum64 struct {
		// digest contains the current state of the Murmur3 64-bit hash
		digest Hash128
		// seed is the Murmur3 64-bit specific seed value
		seed uint64
	}
)

// New64 returns a Murmur3 64-bit hash function that implements the hash.Hash64
// interface with a default seed value of zero.
func New64() hash.Hash64 {
	return New64WithSeed(0)
}

// New64WithSeed returns a Murmur3 64-bit hash function that implements the
// hash.Hash64 interface with a specified seed value.
func New64WithSeed(seed uint64) hash.Hash64 {
	return &sum64{digest: New128WithSeed(seed), seed: seed}
}

// Write implements the io.Writer interface required by the hash.Hash64
// interface. It is the means by which the hash function digests data.
func (s *sum64) Write(data []byte) (int, error) {
	return s.digest.Write(data)
}

// Sum64 implements the hash.Hash64 interface. It returns the Murmur3 64-bit
// hash of the internal state.
func (s *sum64) Sum64() uint64 {
	return s.digest.Sum128()[0]
}

// Sum implements the hash.Hash64 interface. It appends the current hash to b
// and returns the resulting slice. It does not change the underlying
// Murmurhash3 64-bit hash state.
func (s *sum64) Sum(in []byte) []byte {
	h := s.digest.Sum128()[0]

	return append(in,
		byte(h>>56),
		byte(h>>48),
		byte(h>>40),
		byte(h>>32),
		byte(h>>24),
		byte(h>>16),
		byte(h>>8),
		byte(h),
	)
}

// Reset implements the hash.Hash interface. Murmurhash3 128-bit initial state
// is two 64 bit unsigned integers where the initial values are the seed.
func (s *sum64) Reset() {
	s.digest.Reset()
}

// Size implements the hash.Hash64 interface. Murmurhash3 64-bit returns an
// eight byte digest sum.
func (s *sum64) Size() int {
	return 8
}

// BlockSize implements the hash.Hash64 interface.
func (s *sum64) BlockSize() int {
	return 1
}
