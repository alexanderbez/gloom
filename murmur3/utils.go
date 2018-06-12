package murmur3

// rotl64 performs a 64-bit left circular shift count times on n.
func rotl64(n uint64, count uint) uint64 {
	return (n << count) | (n >> (64 - count))
}
