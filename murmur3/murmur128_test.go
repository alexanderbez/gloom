package murmur3

import "testing"

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

		if s.Sum128() != tc.hash {
			t.Errorf("unexpected result: expected: %v, got: %v", tc.hash, s.Sum128())
		}
	}
}
