package wavelettree

import (
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBitvector(t *testing.T) {
	vec := makeBitVector(13)
	vec.Set(3, 0, 4)
	result := vec.Get(3, 0)
	require.Equal(t, uint8(4), result, "Get = previously set")

	vec.Set(1, 3, 1)
	result = vec.Get(1, 3)
	require.Equal(t, uint8(1), result, "Get = previously set, different size")

	vec.Set(4, 7, 15)
	result = vec.Get(4, 7)
	require.Equal(t, uint8(15), result, "Get = previously set, crossing byte")

	vec.Set(3, 1, 6)
	result = vec.Get(3, 1)
	require.Equal(t, uint8(6), result, "Overriding prior")

	vec = vec.Append(3, 7)
	require.Equal(t, uint64(16), vec.bitlength, "Length should = 16")
	result = vec.Get(3, 13)
	require.Equal(t, uint8(7), result, "Get = appended")
}

func TestCalcBlocksize(t *testing.T) {
	// test first 1000 integers
	for n := range 1000 {
		expected := uint8(math.Round(math.Log2(float64(n)) / 2))
		require.Equal(t, expected, calcBlocksize(uint64(n)))
	}
	// test 1000 random integers
	for range 1000 {
		n := rand.Uint64()
		expected := uint8(math.Round(math.Log2(float64(n)) / 2))
		require.Equal(t, expected, calcBlocksize(n))
	}
}

func TestNewRRR(t *testing.T) {
	out := NewRRR(makeBitVector(1000))
	t.Log(out.blockSize, out.classSize, out.offsetSize, out.encoded.bitlength)
}
