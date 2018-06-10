package murmur32

import "hash"

const (
	c1 uint32 = 0xcc9e2d51
	c2 uint32 = 0x1b873593
	r1 uint32 = 15
	r2 uint32 = 13
	m  uint32 = 5
	n  uint32 = 0xe6546b64
)

type (
	// Murmur32 implements the non-cryptographic Murmur3 32-bit hash function
	// suitable for general hash-based lookup.
	Murmur32 struct {
		// sum32 reflects the current state of the hash
		sum32 uint32
		// seed is the Murmur3 32-bit specific seed value
		seed uint32
	}
)

// New32 returns a Murmur3 32-bit hash function that implements the hash.Hash32
// interface with a default seed value of zero.
func New32() hash.Hash32 {
	return New32WithSeed(0)
}

// New32WithSeed returns a Murmur3 32-bit hash function that implements the
// hash.Hash32 interface with a specified seed value.
func New32WithSeed(seed uint32) hash.Hash32 {
	return &Murmur32{sum32: seed, seed: seed}
}

// TWrite implements the io.Writer interface required by the hash.Hash32
// interface. It the means by which the hash function digests data.
func (mm *Murmur32) Write(data []byte) (int, error) {
	remainingIdx := 0
	for i := 0; i+4 <= len(data); i += 4 {
		remainingIdx = i + 4

		// Construct the four byte chunk as a uint32
		a := uint32(data[i])
		b := uint32(data[i+1]) << 8
		c := uint32(data[i+2]) << 16
		d := uint32(data[i+3]) << 24
		k := a | b | c | d

		k *= c1
		k = (k << r1) | (k >> (32 - r1))
		k *= c2

		mm.sum32 ^= k
		mm.sum32 = (mm.sum32 << r2) | (mm.sum32 >> (32 - r2))
		mm.sum32 *= m
		mm.sum32 += n
	}

	// Check if there are any remaining bytes to consume. If so, there can only
	// be one, two, or three bytes remaining.
	if remainingLen := len(data) - remainingIdx; remainingLen > 0 {
		var remainingBytes uint32

		switch remainingLen {
		case 3:
			a := uint32(data[remainingIdx])
			b := uint32(data[remainingIdx+1]) << 8
			c := uint32(data[remainingIdx+2]) << 16
			remainingBytes = a | b | c
		case 2:
			a := uint32(data[remainingIdx])
			b := uint32(data[remainingIdx+1]) << 8
			remainingBytes = a | b
		case 1:
			remainingBytes = uint32(data[remainingIdx])
		}

		remainingBytes *= c1
		remainingBytes = (remainingBytes << r1) | (remainingBytes >> (32 - r1))
		remainingBytes *= c2
		mm.sum32 ^= remainingBytes
	}

	mm.sum32 ^= uint32(len(data))
	mm.sum32 ^= mm.sum32 >> 16
	mm.sum32 *= 0x85ebca6b
	mm.sum32 ^= mm.sum32 >> 13
	mm.sum32 *= 0xc2b2ae35
	mm.sum32 ^= mm.sum32 >> 16

	return len(data), nil
}

// Sum implements the hash.Hash32 interface. It appends the current hash to b
// and returns the resulting slice. It does not change the underlying
// Murmurhash3 32 bit hash state.
func (mm *Murmur32) Sum(in []byte) []byte {
	h := mm.Sum32()
	return append(in, byte(h>>24), byte(h>>16), byte(h>>8), byte(h))
}

// Sum32 implements the hash.Hash32 interface. It returns the Murmurhash3 32
// bit internal hash state.
func (mm *Murmur32) Sum32() uint32 {
	return mm.sum32
}

// Reset implements the hash.Hash32 interface. Murmurhash3 32 bit initial state
// is the seed.
func (mm *Murmur32) Reset() {
	mm.sum32 = mm.seed
}

// Size implements the hash.Hash32 interface. Murmurhash3 32 bit returns a four
// byte digest sum.
func (mm *Murmur32) Size() int {
	return 4
}

// BlockSize implements the hash.Hash32 interface. Murmurhash3 32 bit works in
// chunks of four bytes.
func (mm *Murmur32) BlockSize() int {
	return 4
}
