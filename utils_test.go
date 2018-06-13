package gloom

import "testing"

func TestOptimalBitVectorSize(t *testing.T) {
	testCases := []struct {
		n uint64
		p float64
		e uint64
	}{
		{n: 10, p: 0.04, e: 67},
		{n: 5000, p: 0.01, e: 47926},
		{n: 100000, p: 0.01, e: 958506},
	}

	for _, tc := range testCases {
		m := optimalBitVectorSize(tc.n, tc.p)

		if m != tc.e {
			t.Errorf("unexpected result: expected: %v, got: %v", tc.e, m)
		}
	}
}

func TestOptimalNumHash(t *testing.T) {
	testCases := []struct {
		m, n uint64
		e    uint64
	}{
		{m: 67, n: 10, e: 5},
		{m: 47926, n: 5000, e: 7},
		{m: 958506, n: 100000, e: 7},
	}

	for _, tc := range testCases {
		k := optimalNumHash(tc.m, tc.n)

		if k != tc.e {
			t.Errorf("unexpected result: expected: %v, got: %v", tc.e, k)
		}
	}
}
