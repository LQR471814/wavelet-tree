package wavelettree

import (
	"testing"
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
