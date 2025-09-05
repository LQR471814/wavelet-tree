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
	// 65k
	n := ^uint16(0)

	pool := sync.Pool{
		New: func() any {
			vec := NewBitVector(uint64(n))
			return &vec
		},
	}

	// initialize with 1024 objects
	for range 1024 {
		pool.Put(pool.New())
	}

	randIndex := func(rndm *rand.Rand, bitsize uint8) uint64 {
		upper := n - uint16(bitsize)
		return uint64(rndm.Intn(int(upper)))
	}

	f.Add(int64(4))
	f.Fuzz(func(t *testing.T, seed int64) {
		rndm := rand.New(rand.NewSource(seed))

		vec := pool.Get().(*BitVector)
		defer pool.Put(vec)

		bytesize := rndm.Intn(4)

		switch bytesize {
		case 0:
			value := uint8(rndm.Uint64() % uint64(^uint8(0)))
			bitsize := uint8(bits.Len8(value))
			if bitsize == 0 {
				bitsize = 1
			}
			rndm.Uint32()
			i := randIndex(rndm, bitsize)
			vec.Set8(bitsize, i, value)

			result := vec.Get8(bitsize, i)
			require.Equal(t, value, result, "Set-Get not equal (uint8)")
		case 1:
			value := uint16(rndm.Uint64() % uint64(^uint16(0)))
			bitsize := uint8(bits.Len16(value))
			if bitsize == 0 {
				bitsize = 1
			}
			i := randIndex(rndm, bitsize)
			vec.Set16(bitsize, i, value)

			result := vec.Get16(bitsize, i)
			require.Equal(t, value, result, "Set-Get not equal (uint16)")
		case 2:
			value := uint32(rndm.Uint64() % uint64(^uint32(0)))
			bitsize := uint8(bits.Len32(value))
			if bitsize == 0 {
				bitsize = 1
			}
			i := randIndex(rndm, bitsize)
			vec.Set32(bitsize, i, value)

			result := vec.Get32(bitsize, i)
			require.Equal(t, value, result, "Set-Get not equal (uint32)")
		case 3:
			value := uint64(rndm.Uint64() % uint64(^uint64(0)))
			bitsize := uint8(bits.Len64(value))
			if bitsize == 0 {
				bitsize = 1
			}
			i := randIndex(rndm, bitsize)
			vec.Set64(bitsize, i, value)

			result := vec.Get64(bitsize, i)
			require.Equal(t, value, result, "Set-Get not equal (uint64)")
		}
	})
}

func TestBitvector(t *testing.T) {
	{
		vec := NewBitVector(13)
		vec.Set8(3, 0, 4)
		result := vec.Get8(3, 0)
		require.Equal(t, uint8(4), result, "Get = previously set")

		vec.Set8(1, 3, 1)
		result = vec.Get8(1, 3)
		require.Equal(t, uint8(1), result, "Get = previously set, different size")

		vec.Set8(4, 7, 15)
		result = vec.Get8(4, 7)
		require.Equal(t, uint8(15), result, "Get = previously set, crossing byte")

		vec.Set8(3, 1, 6)
		result = vec.Get8(3, 1)
		require.Equal(t, uint8(6), result, "Overriding prior")

		vec = vec.Append8(3, 7)
		require.Equal(t, uint64(16), vec.bitlength, "Length should = 16")
		result = vec.Get8(3, 13)
		require.Equal(t, uint8(7), result, "Get = appended")
	}
}
