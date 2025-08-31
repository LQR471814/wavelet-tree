package wavelettree

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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

func TestNewRRR(t *testing.T) {
	out := NewRRR(NewBitVector(10000), RRROptions{})
	t.Log(
		"block size:",
		out.blockSize,
	)
	t.Log(
		"superblock size:",
		out.superblockSize,
	)
	t.Log(
		"size(class):",
		out.classFieldSize,
	)
	t.Log(
		"size(offset)",
		out.offsetFieldSize,
	)
	t.Log(
		"size(cumulative rank)",
		out.cumulativeRankFieldSize,
	)
	t.Log(
		"size(all):",
		out.bits.bitlength,
	)
}
