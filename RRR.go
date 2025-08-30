package wavelettree

import (
	"slices"
)

/*
RRR description:

divide bitvector into fixed-size blocks
usually b (size) = log(n) / 2 bits (this is the optimal value)

for each block store:
  - the # of 1's (class)
  - the position of 1's (offset)
      - the offset is the index of the current block in the list of possible
      combinations of patterns with the given block size and class
      - ex. 1 0 1 1
      - 4 bits has C(4,3)

superblocks which are blocks of blocks can be used to accelerate rank
queries.
  - each superblock contains the total number of 1s of the blocks inside it
  - so if you were to compute `rank(i)`, instead of having to go through
  each of the blocks up to `i`, you would only need to sum up all the
  "superblocks" until `i`, then go through the block which contains `i`

for efficient operations, we'll want to consider the CPU's word size. (for
64 bit cpus, this would be 64 bits) This means that the practical upper limit
for block size b is defined as:

b = min(log(n)/2, 64)

this also indicates that our max size bitvector n is given by:
log(n)/2 = 64
n = 10^128

which should frankly be plenty, so we don't need to worry about reaching the
maximum block size.
*/

/*
Blocks:

block configurations can vary based on block size:
- b <= 8
	- 3 bits for class
- b <= 4
	- 2 bits for class

`log_2(C(b, class))` bits for offset.
so `C(4, 3)` for b=4 and class=3 would yield `log_2(4)` which is `2`

the most amount of memory an offset field could take is 7 bits
that is the result of `ceil(log_2(C(8, 4)))`.

block is encoded as follows:
- class (some bits)
- offset (some bits)
*/

// RRR enables near O(1) calculations of bit rank(i) and other operations.
type RRR struct {
	encoded BitVector
	// blockSize is the number of bits in a block (value from [1, 64])
	blockSize uint8
	// classFieldSize (number of bits required to store the number of 1s for each
	// block, max: # of bits in the block)
	classFieldSize uint8
	// offsetFieldSize (number of bits required to store the offset for each block,
	// max: C(n, n/2) + 1)
	offsetFieldSize uint8
}

// maps class -> possible offset combinations
var offset_lookup_uint8 [][]uint8
var offset_lookup_uint16 [][]uint16
var offset_lookup_uint32 [][]uint32
var offset_lookup_uint64 [][]uint64

func computeOffset() {
}

func computeOffsetLookup[T uint8 | uint16 | uint32 | uint64](out *[][]T, bytesize uint8) {
	slice := make([][]T, bytesize+1)
	for value := T(0); value <= ^T(0); value++ {
		c := countbits(bytesize, value)
		slice[c] = append(slice[c], value)
	}
	*out = slice
}

func init() {
	computeOffsetLookup(&offset_lookup_uint8, 8)
	computeOffsetLookup(&offset_lookup_uint16, 16)
	computeOffsetLookup(&offset_lookup_uint32, 32)
	computeOffsetLookup(&offset_lookup_uint64, 64)
}

func getBlockValues(blockSize uint8, i uint64, bits BitVector) (class, offset uint8) {
	switch {
	case blockSize <= 8:
		content := bits.Get8(blockSize, i)
		class = countbits(8, content)
		offset = uint8(slices.Index(offset_lookup_uint8[class], content))
		return
	case blockSize <= 16:
		content := bits.Get16(blockSize, i)
		class = countbits(16, content)
		offset = uint8(slices.Index(offset_lookup_uint16[class], content))
		return
	case blockSize <= 32:
		content := bits.Get32(blockSize, i)
		class = countbits(32, content)
		offset = uint8(slices.Index(offset_lookup_uint32[class], content))
		return
	case blockSize <= 64:
		content := bits.Get64(blockSize, i)
		class = countbits(64, content)
		offset = uint8(slices.Index(offset_lookup_uint64[class], content))
		return
	}
	panic("exceeded max block length 64!")
}

func NewRRR(bits BitVector) (out RRR) {
	n := bits.Length()

	blocksize := floorLog2(n)
	blocksize >>= 1
	out.blockSize = blocksize
	out.classFieldSize = floorLog2(out.blockSize)

	maxOffset := choose(uint64(out.blockSize), uint64(out.blockSize)>>1)
	out.offsetFieldSize = floorLog2(maxOffset)

	blocks := n/uint64(out.blockSize) + 1
	blockSize := out.classFieldSize + out.offsetFieldSize
	totalSize := blocks * uint64(blockSize)
	out.encoded = NewBitVector(totalSize)

	for i := range blocks {
		bitIdx := i * uint64(blockSize)

		class, offset := getBlockValues(out.blockSize, bitIdx, bits)
		out.encoded.Set8(out.classFieldSize, bitIdx, class)
		out.encoded.Set8(out.offsetFieldSize, bitIdx+uint64(out.classFieldSize), offset)
	}

	return
}

// func (r RRR) Rank() {
//
// }
