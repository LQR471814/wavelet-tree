package wavelettree

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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

func TestRank(t *testing.T) {
	// either: 001, 010, 100
	res := rank[uint8](3, 0b010)
	require.Equal(t, uint64(1), res)

	// either: 011, 101, 110
	res = rank[uint8](3, 0b110)
	require.Equal(t, uint64(2), res)

	// only: 000
	res = rank[uint8](3, 0b000)
	require.Equal(t, uint64(0), res)

	// only: 111
	res = rank[uint8](3, 0b111)
	require.Equal(t, uint64(0), res)
}

func TestUnrank(t *testing.T) {
	require.Equal(t, uint8(0b010), unrank[uint8](1, 1))
	require.Equal(t, uint8(0b110), unrank[uint8](2, 2))
	require.Equal(t, uint8(0b000), unrank[uint8](0, 0))
	require.Equal(t, uint8(0b111), unrank[uint8](3, 0))
}
