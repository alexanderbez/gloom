package murmur3

import "testing"

func TestSum64(t *testing.T) {
	testCases := []struct {
		seed uint64
		data []byte
		hash uint64
	}{
		{
			data: nil,
			hash: 0x0,
		},
		{
			data: []byte{},
			hash: 0x0,
		},
		{
			data: []byte("Hello, World!!!!"),
			hash: 0xb57df55e4edee585,
		},
		{
			seed: 0xfa,
			data: []byte("Hello, World!!!!"),
			hash: 0xb02fba6b1d629cae,
		},
		{
			data: []byte("Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa."),
			hash: 0xaf11090ad904f11a,
		},
		{
			seed: 0xfa,
			data: []byte("Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa."),
			hash: 0x2d36aa481cf715ec,
		},
		{
			data: []byte("@@@##!&^#%$!+_][;//"),
			hash: 0xae66907cdc6d6934,
		},
		{
			seed: 0xfa,
			data: []byte("@@@##!&^#%$!+_][;//"),
			hash: 0x9ec3da2f2c2441b5,
		},
	}

	for _, tc := range testCases {
		s := New64WithSeed(tc.seed)
		s.Write(tc.data)

		if s.Sum64() != tc.hash {
			t.Errorf("unexpected result: expected: %v, got: %v", tc.hash, s.Sum64())
		}
	}
}
