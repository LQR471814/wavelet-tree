package wavelettree

import (
	"math/bits"
	"math/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// Invariants:
// - no crashes of any kind
// - set(size, i, value) -> get(size, i, value) should be equal
// - set(size, i, value) -> get(size, j \in !(i, i+size), value) should not change
// - set(size, i, value) 2x -> get(size, i, value) should be equal

func FuzzBitVector(f *testing.F) {
	// ~4Gbits
	n := ^uint32(0)

	pool := sync.Pool{
		New: func() any {
			return NewBitVector(uint64(n))
		},
	}

	randIndex := func(rndm *rand.Rand, bitsize uint8) uint64 {
		upper := n - uint32(bitsize)
		r := rndm.Uint32()
		for r > upper {
			r = rndm.Uint32()
		}
		return uint64(r)
	}

	f.Add(int64(4))
	f.Fuzz(func(t *testing.T, seed int64) {
		rndm := rand.New(rand.NewSource(seed))

		vec := pool.Get().(BitVector)
		defer pool.Put(vec)

		bytesize := rndm.Intn(4)

		switch bytesize {
		case 0:
			value := uint8(rndm.Uint64() % uint64(^uint8(0)))
			bitsize := uint8(bits.Len8(value))
			rndm.Uint32()
			i := randIndex(rndm, bitsize)
			vec.Set8(bitsize, i, value)

			result := vec.Get8(bitsize, i)
			require.Equal(t, value, result, "Set-Get not equal (uint8)")
		case 1:
			value := uint16(rndm.Uint64() % uint64(^uint16(0)))
			bitsize := uint8(bits.Len16(value))
			i := randIndex(rndm, bitsize)
			vec.Set16(bitsize, i, value)

			result := vec.Get16(bitsize, i)
			require.Equal(t, value, result, "Set-Get not equal (uint16)")
		case 2:
			value := uint32(rndm.Uint64() % uint64(^uint32(0)))
			bitsize := uint8(bits.Len32(value))
			i := randIndex(rndm, bitsize)
			vec.Set32(bitsize, i, value)

			result := vec.Get32(bitsize, i)
			require.Equal(t, value, result, "Set-Get not equal (uint32)")
		case 3:
			value := uint64(rndm.Uint64() % uint64(^uint64(0)))
			bitsize := uint8(bits.Len64(value))
			i := randIndex(rndm, bitsize)
			vec.Set64(bitsize, i, value)

			result := vec.Get64(bitsize, i)
			require.Equal(t, value, result, "Set-Get not equal (uint64)")
		}
	})
}
