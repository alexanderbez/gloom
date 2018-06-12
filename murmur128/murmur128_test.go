package murmur128

import "testing"

func TestNew32(t *testing.T) {
	testCases := []struct {
		data []byte
		hash uint32
	}{
		{
			data: nil,
			hash: 0,
		},
		{
			data: []byte{},
			hash: 0,
		},
		{
			data: []byte("Hello, world!"),
			hash: 3224780355,
		},
		{
			data: []byte("Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa."),
			hash: 2390704954,
		},
		{
			data: []byte("@@@##!&^#%$!+_][;//"),
			hash: 1089243445,
		},
	}

	mm := New32()

	for _, tc := range testCases {
		mm.Reset()
		mm.Write(tc.data)

		if mm.Sum32() != tc.hash {
			t.Errorf("unexpected result: expected: %v, got: %v", tc.hash, mm.Sum32())
		}
	}
}

func TestNew32WithSeed(t *testing.T) {
	testCases := []struct {
		data []byte
		hash uint32
	}{
		{
			data: nil,
			hash: 454360506,
		},
		{
			data: []byte{},
			hash: 454360506,
		},
		{
			data: []byte("Hello, world!"),
			hash: 155145377,
		},
		{
			data: []byte("Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa."),
			hash: 35112277,
		},
		{
			data: []byte("@@@##!&^#%$!+_][;//"),
			hash: 2493164017,
		},
	}

	largePrime := uint32(15487469)
	mm := New32WithSeed(largePrime)

	for _, tc := range testCases {
		mm.Reset()
		mm.Write(tc.data)

		if mm.Sum32() != tc.hash {
			t.Errorf("unexpected result: expected: %v, got: %v", tc.hash, mm.Sum32())
		}
	}
}