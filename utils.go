package wavelettree

type integers interface {
	int | uint |
		int8 | uint8 |
		int16 | uint16 |
		int32 | uint32 |
		int64 | uint64
}

// floor(log_2(n)) efficiently
func floorLog2[T integers](n T) (out uint8) {
	var zero T
	for n > zero {
		n >>= 1
		out++
	}
	return
}

// this is realistically as efficient as we'll get for choose with small (sub
// multi-million) n
func choose(n, k uint64) (result uint64) {
	if k > n {
		return 0
	}
	if k == 0 || k == n {
		return 1
	}
	if k > n-k {
		k = n - k
	}
	result = 1
	for i := uint64(1); i <= k; i++ {
		result = result * (n - i + 1) / i
	}
	return result
}
