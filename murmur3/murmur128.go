package murmur3

import (
	"hash"
)

const (
	c1 uint64 = 0x87c37b91114253d5
	c2 uint64 = 0x4cf5ad432745937f
)

type (
	// sum128 implements the non-cryptographic Murmur3 128-bit hash function
	// suitable for general hash-based lookup. The internal state is composed
	// of two 64-bit unsigned integers where both have an initial seed value.
	sum128 struct {
		// digest contains the current state of the Murmur3 128-bit hash
		digest [2]uint64
		// seed is the Murmur3 128-bit specific seed value
		seed uint64
	}

	// Hash128 reflects a 128-bit hash interface
	Hash128 interface {
		hash.Hash
		Sum128() [2]uint64
	}
)

// New128 returns a Murmur3 128-bit hash function that implements the Hash128
// interface with a default seed value of zero.
func New128() Hash128 {
	return New128WithSeed(0)
}

// New128WithSeed returns a Murmur3 128-bit hash function that implements the
// Hash128 interface with a specified seed value.
func New128WithSeed(seed uint64) Hash128 {
	return &sum128{digest: [2]uint64{seed, seed}, seed: seed}
}

// Sum128 returns the Murmur3 128-bit sum of the hash's state.
func (s *sum128) Sum128() [2]uint64 {
	return s.digest
}

// Write implements the io.Writer interface required by the hash.Hash
// interface. It is the means by which the hash function digests data.
func (s *sum128) Write(data []byte) (int, error) {
	tail := s.blockInterMix(data)
	s.blockInterMixTail(tail)
	s.finalizeMix(data)

	return len(data), nil
}

// blockInterMix performs inter-block mixing on a Murmur3 128-bit digest and
// the provided data such that the blocks digested are 128-bits wide. Any
// remaining non-digested bytes are returned for further inter-block mixing.
func (s *sum128) blockInterMix(data []byte) []byte {
	tailStart := 0
	for i := 0; i+s.Size() <= len(data); i += s.Size() {
		var k1, k2 uint64

		tailStart += s.Size()

		// Get the two respective 64-bit chunks from the 128-bit chunk as
		// 64-bit unsigned integers.
		x := data[i : i+s.Size()/2]
		y := data[i+s.Size()/2:]

		for j := 0; j < s.Size()/2; j++ {
			k1 |= uint64(x[j]) << uint(j*8)
			k2 |= uint64(y[j]) << uint(j*8)
		}

		k1 *= c1
		k1 = rotl64(k1, 31)
		k1 *= c2

		s.digest[0] ^= k1
		s.digest[0] = rotl64(s.digest[0], 27)
		s.digest[0] += s.digest[1]
		s.digest[0] = s.digest[0]*5 + 0x52dce729

		k2 *= c2
		k2 = rotl64(k2, 33)
		k2 *= c1

		s.digest[1] ^= k2
		s.digest[1] = rotl64(s.digest[1], 31)
		s.digest[1] += s.digest[0]
		s.digest[1] = s.digest[1]*5 + 0x38495ab5
	}

	return data[tailStart:]
}

// blockInterMixTail performs inter-block mixing on a Murmur3 128-bit digest
// and the provided data such that the blocks digested are the remainder after
// all 128-bit wide blocks have been digested.
func (s *sum128) blockInterMixTail(tail []byte) {
	var k1, k2 uint64

	switch len(tail) & 15 {
	case 15:
		k2 ^= uint64(tail[14]) << 48
		fallthrough
	case 14:
		k2 ^= uint64(tail[13]) << 40
		fallthrough
	case 13:
		k2 ^= uint64(tail[12]) << 32
		fallthrough
	case 12:
		k2 ^= uint64(tail[11]) << 24
		fallthrough
	case 11:
		k2 ^= uint64(tail[10]) << 16
		fallthrough
	case 10:
		k2 ^= uint64(tail[9]) << 8
		fallthrough
	case 9:
		k2 ^= uint64(tail[8]) << 0
		k2 *= c2
		k2 = rotl64(k2, 33)
		k2 *= c1

		s.digest[1] ^= k2
		fallthrough
	case 8:
		k1 ^= uint64(tail[7]) << 56
		fallthrough
	case 7:
		k1 ^= uint64(tail[6]) << 48
		fallthrough
	case 6:
		k1 ^= uint64(tail[5]) << 40
		fallthrough
	case 5:
		k1 ^= uint64(tail[4]) << 32
		fallthrough
	case 4:
		k1 ^= uint64(tail[3]) << 24
		fallthrough
	case 3:
		k1 ^= uint64(tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint64(tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint64(tail[0]) << 0
		k1 *= c1
		k1 = rotl64(k1, 31)
		k1 *= c2

		s.digest[0] ^= k1
	}
}

// finalizeMix performs a finalization mix for the Murmur3 128-bit digest.
func (s *sum128) finalizeMix(data []byte) {
	s.digest[0] ^= uint64(len(data))
	s.digest[1] ^= uint64(len(data))

	s.digest[0] += s.digest[1]
	s.digest[1] += s.digest[0]

	s.digest[0] ^= s.digest[0] >> 33
	s.digest[0] *= 0xff51afd7ed558ccd
	s.digest[0] ^= s.digest[0] >> 33
	s.digest[0] *= 0xc4ceb9fe1a85ec53
	s.digest[0] ^= s.digest[0] >> 33

	s.digest[1] ^= s.digest[1] >> 33
	s.digest[1] *= 0xff51afd7ed558ccd
	s.digest[1] ^= s.digest[1] >> 33
	s.digest[1] *= 0xc4ceb9fe1a85ec53
	s.digest[1] ^= s.digest[1] >> 33

	s.digest[0] += s.digest[1]
	s.digest[1] += s.digest[0]

}

// Sum implements the hash.Hash interface. It appends the current hash to b
// and returns the resulting slice. It does not change the underlying
// Murmurhash3 128-bit hash state.
func (s *sum128) Sum(in []byte) []byte {
	// Append in MSB order chunks starting the most significant (first) 64-bit
	// unsigned integer.
	return append(in,
		byte(s.digest[0]>>56),
		byte(s.digest[0]>>48),
		byte(s.digest[0]>>40),
		byte(s.digest[0]>>32),
		byte(s.digest[0]>>24),
		byte(s.digest[0]>>16),
		byte(s.digest[0]>>8),
		byte(s.digest[0]),
		byte(s.digest[1]>>56),
		byte(s.digest[1]>>48),
		byte(s.digest[1]>>40),
		byte(s.digest[1]>>32),
		byte(s.digest[1]>>24),
		byte(s.digest[1]>>16),
		byte(s.digest[1]>>8),
		byte(s.digest[1]),
	)

}

// Reset implements the hash.Hash interface. Murmurhash3 128-bit initial state
// is two 64 bit unsigned integers where the initial values are the seed.
func (s *sum128) Reset() {
	s.digest[0] = s.seed
	s.digest[1] = s.seed
}

// Size implements the hash.Hash interface. Murmurhash3 128-bit returns a 16
// byte digest sum.
func (s *sum128) Size() int {
	return 16
}

// BlockSize implements the hash.Hash interface.
func (s *sum128) BlockSize() int {
	return 1
}
