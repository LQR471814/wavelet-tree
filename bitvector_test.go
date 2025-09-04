package wavelettree

import (
	"fmt"
	"math/rand"
	"testing"
)

// Invariants:
// - no crashes of any kind
// - set(size, i, value) -> get(size, i, value) should be equal
// - set(size, i, value) -> get(size, j \in !(i, i+size), value) should not change
// - set(size, i, value) 2x -> get(size, i, value) should be equal

func randUint64n(rndm *rand.Rand, n uint64) uint64 {
	v := rndm.Uint64()
	for v > n {
		v = rndm.Uint64()
	}
	return v
}

func FuzzBitVector(f *testing.F) {
	// n := uint64(1) << 63

	// vec := NewBitVector(n)
	// lock := sync.Mutex{}

	fmt.Println("here!")
	// rndm := rand.New(rand.NewSource(2))

	// f.Add(int64(4))
	// f.Fuzz(func(t *testing.T, seed int64) {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	//
	// 	bytesize := rndm.Intn(4)
	//
	// 	switch bytesize {
	// 	case 0:
	// 		value := uint8(rndm.Uint64() % uint64(^uint8(0)))
	// 		bitsize := uint8(bits.Len8(value))
	// 		i := randUint64n(rndm, n-uint64(bitsize))
	// 		vec.Set8(bitsize, i, value)
	//
	// 		result := vec.Get8(bitsize, i)
	// 		require.Equal(t, value, result, "Set-Get not equal")
	// 	case 1:
	// 		value := uint16(rndm.Uint64() % uint64(^uint16(0)))
	// 		bitsize := uint8(bits.Len16(value))
	// 		i := randUint64n(rndm, n-uint64(bitsize))
	// 		vec.Set16(bitsize, i, value)
	//
	// 		result := vec.Get16(bitsize, i)
	// 		require.Equal(t, value, result, "Set-Get not equal")
	// 	case 2:
	// 		value := uint32(rndm.Uint64() % uint64(^uint32(0)))
	// 		bitsize := uint8(bits.Len32(value))
	// 		i := randUint64n(rndm, n-uint64(bitsize))
	// 		vec.Set32(bitsize, i, value)
	//
	// 		result := vec.Get32(bitsize, i)
	// 		require.Equal(t, value, result, "Set-Get not equal")
	// 	case 3:
	// 		value := uint64(rndm.Uint64() % uint64(^uint64(0)))
	// 		bitsize := uint8(bits.Len64(value))
	// 		i := randUint64n(rndm, n-uint64(bitsize))
	// 		vec.Set64(bitsize, i, value)
	//
	// 		result := vec.Get64(bitsize, i)
	// 		require.Equal(t, value, result, "Set-Get not equal")
	// 	}
	// })
}
