package gloom

import (
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestBloomFilter(t *testing.T) {
	items := map[string]struct{}{}

	var (
		n uint64 = 1000
		r uint64 = 500
		i uint64
	)

	token := make([]byte, 20)
	for ; i < n; i++ {
		rand.Read(token)
		items[string(token)] = struct{}{}
	}

	bf, err := NewBloomFilter(uint64(len(items)), DefaultFalsePosProb)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	for item := range items {
		if err := bf.Set([]byte(item)); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}

	for item := range items {
		ok, err := bf.Includes([]byte(item))

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !ok {
			t.Errorf("unexpected result: expected true for %s", item)
		}
	}

	i = 0
	for ; i < r; i++ {
		rand.Read(token)

		ok, err := bf.Includes(token)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if _, inSet := items[string(token)]; !ok && inSet {
			t.Errorf("unexpected result: false negative for %s", string(token))
		}
	}
}
