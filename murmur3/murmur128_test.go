package murmur3

import (
	"testing"

	mm3 "github.com/spaolacci/murmur3"
)

func TestSum128(t *testing.T) {
	testCases := []struct {
		seed uint64
		data []byte
		hash [2]uint64
	}{
		{
			data: nil,
			hash: [2]uint64{0x0, 0x0},
		},
		{
			data: []byte{},
			hash: [2]uint64{0x0, 0x0},
		},
		{
			data: []byte("Hello, World!!!!"),
			hash: [2]uint64{0xb57df55e4edee585, 0x7fbabd101b969fb2},
		},
		{
			seed: 0xfa,
			data: []byte("Hello, World!!!!"),
			hash: [2]uint64{0xb02fba6b1d629cae, 0xab907abfb23bfcfc},
		},
		{
			data: []byte("Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa."),
			hash: [2]uint64{0xaf11090ad904f11a, 0x52b5309456f0ad38},
		},
		{
			seed: 0xfa,
			data: []byte("Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa."),
			hash: [2]uint64{0x2d36aa481cf715ec, 0x15f04536764cf671},
		},
		{
			data: []byte("@@@##!&^#%$!+_][;//"),
			hash: [2]uint64{0xae66907cdc6d6934, 0xf8e73c715a15b592},
		},
		{
			seed: 0xfa,
			data: []byte("@@@##!&^#%$!+_][;//"),
			hash: [2]uint64{0x9ec3da2f2c2441b5, 0xfe8d97a6d3920f83},
		},
	}

	for _, tc := range testCases {
		s := New128WithSeed(tc.seed)
		s.Write(tc.data)

		h1, h2 := mm3.Sum128WithSeed(tc.data, uint32(tc.seed))

		if h1 != tc.hash[0] {
			t.Errorf("unexpected result: expected: %v, got: %v", tc.hash[0], h1)
		}

		if h2 != tc.hash[1] {
			t.Errorf("unexpected result: expected: %v, got: %v", tc.hash[1], h2)
		}

		if s.Sum128() != tc.hash {
			t.Errorf("unexpected result: expected: %v, got: %v", tc.hash, s.Sum128())
		}
	}
}
