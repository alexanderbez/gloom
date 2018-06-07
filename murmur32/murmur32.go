package gloom

import "hash"

// type Hash interface {
// 	// Write (via the embedded io.Writer interface) adds more data to the running hash.
// 	// It never returns an error.
// 	io.Writer

// 	Sum(b []byte) []byte

type (
	murmur32 struct {
		sum32 uint32
		seed  uint32
	}
)

func New32() hash.Hash32 {
	return New32WithSeed(0)
}

func New32WithSeed(seed uint32) hash.Hash32 {
	return &murmur32{sum32: seed, seed: seed}
}

// Sum implements the hash.Hash32 interface. It appends the current hash to b
// and returns the resulting slice. It does not change the underlying
// Murmurhash3 32 bit hash state.
func (m *murmur32) Sum(in []byte) []byte {

}

// Sum32 implements the hash.Hash32 interface. It returns the Murmurhash3 32
// bit internal hash state.
func (m *murmur32) Sum32() uint32 {
	return m.sum32
}

// Reset implements the hash.Hash32 interface. Murmurhash3 32 bit initial state
// is the seed.
func (m *murmur32) Reset() {
	m.sum32 = m.seed
}

// Size implements the hash.Hash32 interface. Murmurhash3 32 bit returns a four
// byte digest sum.
func (m *murmur32) Size() int {
	return 4
}

// BlockSize implements the hash.Hash32 interface. Murmurhash3 32 bit works in
// chunks of four bytes.
func (m *murmur32) BlockSize() int {
	return 4
}
