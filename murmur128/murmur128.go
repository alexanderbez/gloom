package murmur128

import "hash"

const (
	c1 uint64 = 0x87c37b91114253d5
	c2 uint64 = 0x4cf5ad432745937f
)

type (
	// Murmur128 implements the non-cryptographic Murmur3 128-bit hash function
	// suitable for general hash-based lookup. The internal state is composed
	// of two 64-bit unsigned integers where both have an initial seed value.
	Murmur128 struct {
		// sum128 contains the current state of the Murmur3 128-bit hash
		sum128 [2]uint64
		// seed is the Murmur3 128-bit specific seed value
		seed uint64
	}
)

// New128 returns a Murmur3 128-bit hash function that implements the hash.Hash
// interface with a default seed value of zero.
func New128() hash.Hash {
	return New128WithSeed(0)
}

// New128WithSeed returns a Murmur3 128-bit hash function that implements the
// hash.Hash interface with a specified seed value.
func New128WithSeed(seed uint64) hash.Hash {
	return &Murmur128{sum128: [2]uint64{seed, seed}, seed: seed}
}

// Write implements the io.Writer interface required by the hash.Hash
// interface. It is the means by which the hash function digests data.
func (mm *Murmur128) Write(data []byte) (int, error) {
	mm.blockInterMix(data)

	// TODO: Block inter-mix remaining bits

	// TODO: Finalize

	return len(data), nil
}

// TODO: ...
func (mm *Murmur128) blockInterMix(data []byte) {
	for i := 0; i+mm.Size() <= len(data); i += mm.Size() {
		var (
			k1, k2 uint64
		)

		// Get the two respective 64-bit chunks from the 128-bit chunk as
		// 64-bit unsigned integers.
		x := data[i : i+mm.Size()/2]
		y := data[mm.Size()/2:]

		for j := 0; j < mm.Size()/2; j++ {
			k1 |= uint64(x[j]) << uint(j*8)
			k2 |= uint64(y[j]) << uint(j*8)
		}

		k1 *= c1
		k1 = rotl64(k1, 31)
		k1 *= c2

		mm.sum128[0] ^= k1
		mm.sum128[0] = rotl64(mm.sum128[0], 27)
		mm.sum128[0] += mm.sum128[1]
		mm.sum128[0] = mm.sum128[0]*5 + 0x52dce729

		k2 *= c2
		k2 = rotl64(k2, 33)
		k2 *= c1

		mm.sum128[1] ^= k2
		mm.sum128[1] = rotl64(mm.sum128[1], 31)
		mm.sum128[1] += mm.sum128[0]
		mm.sum128[1] = mm.sum128[1]*5 + 0x38495ab5
	}
}

func (mm *Murmur128) blockInterMixRem() {

}

// Sum implements the hash.Hash interface. It appends the current hash to b
// and returns the resulting slice. It does not change the underlying
// Murmurhash3 128-bit hash state.
func (mm *Murmur128) Sum(in []byte) []byte {
	// Append in MSB order chunks starting the most significant (first) 64-bit
	// unsigned integer.
	return append(in) // byte(h>>56),
	// byte(h>>48),
	// byte(h>>40),
	// byte(h>>32),
	// byte(h>>24),
	// byte(h>>16),
	// byte(h>>8),
	// byte(h),

}

// Reset implements the hash.Hash interface. Murmurhash3 128-bit initial state
// is two 64 bit unsigned integers where the initial values are the seed.
func (mm *Murmur128) Reset() {
	mm.sum128[0] = mm.seed
	mm.sum128[1] = mm.seed
}

// Size implements the hash.Hash interface. Murmurhash3 128-bit returns a 16
// byte digest sum.
func (mm *Murmur128) Size() int {
	return 16
}

// BlockSize implements the hash.Hash interface.
func (mm *Murmur128) BlockSize() int {
	return 1
}

// rotl64 performs a 64-bit left circular count shift on n.
func rotl64(n uint64, count uint) uint64 {
	return (n << count) | (n >> (64 - count))
}
