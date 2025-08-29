package wavelettree

import (
	"fmt"
)

/*
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
block configurations can vary based on block size:
- b <= 8
	- 3 bits for class
- b <= 4
	- 2 bits for class

`log_2(C(b, class))` bits for offset.
so `C(4, 3)` for b=4 and class=3 would yield `log_2(4)` which is `2`

the most amount of memory an offset field could take is 7.
that is the result of `ceil(log_2(C(8, 4)))`.
*/

// bitvector is effectively a slice of bits
type bitvector struct {
	bitlength uint64
	bytes     []uint8
}

func makeBitVector(bitlength uint64) bitvector {
	bytelength := bitlength/8 + 1
	vec := bitvector{
		bitlength: bitlength,
		bytes:     make([]uint8, bytelength),
	}
	return vec
}

// Get gets the value of any size 1-8 bit word at any bit index i from the
// bitvector
func (v bitvector) Get(size uint8, i uint64) (result uint8) {
	if size < 1 || size > 8 {
		panic(fmt.Sprintf("invalid size: %d", size))
	}

	byte := i / 8
	bit := uint8(i % 8)

	var mask uint8 = 255 >> (8 - size)

	result = (v.bytes[byte] >> bit) & mask

	// 9 = 8+1 (since overlap threshold for size=1 should be 8)
	var overlapThreshold uint8 = 9 - size

	// amount of bits set in the current byte
	var currentSet = 8 - bit

	// overlap amount is: amount of bits to be found in the next byte after
	// bits are found in current byte
	var overlapAmount uint8 = size - currentSet

	if bit >= overlapThreshold {
		next := v.bytes[byte+1]
		var nextmask uint8 = 255 >> overlapAmount
		var overlap uint8 = next & nextmask
		result = result | (overlap << currentSet)
	}

	return
}

// Set sets the value of any size 1-8 bit word at any bit index i on the
// bitvector
func (v bitvector) Set(size uint8, i uint64, value uint8) {
	if size < 1 || size > 8 {
		panic(fmt.Sprintf("invalid size: %d", size))
	}

	byte := i / 8
	bit := uint8(i % 8)

	var mask uint8 = 255 >> (8 - size)

	value = value & mask

	surrounding := v.bytes[byte] & (^mask)
	v.bytes[byte] = surrounding | (value << bit)

	// 9 = 8+1 (since overlap threshold for size=1 should be 8)
	var overlapThreshold uint8 = 9 - size

	// amount of bits set in the current byte
	var currentSet uint8 = 8 - bit

	// overlap amount is: amount of bits to be found in the next byte after
	// bits are found in current byte
	var overlapAmount uint8 = size - currentSet

	if bit >= overlapThreshold {
		next := v.bytes[byte+1]
		var nextMask uint8 = 255 >> overlapAmount
		nextupSurround := next & (^nextMask)

		// remove the bits in value that have already been set in the current
		// byte and set those bits in the next byte
		v.bytes[byte+1] = nextupSurround | (value >> currentSet)
	}
}

// Length returns the bit length of the bitvector.
func (v bitvector) Length() uint64 {
	return v.bitlength
}

// Append adds a number of bits 1-8 to the bitvector.
func (v bitvector) Append(size, value uint8) bitvector {
	originalEnd := v.bitlength
	v.bitlength += uint64(size)
	byteLength := v.bitlength/8 + 1
	if int(byteLength) > len(v.bytes) {
		v.bytes = append(v.bytes, 0)
	}
	v.Set(size, originalEnd, value)
	return v
}

// block is encoded as follows
// - class
// - offset

type RRR struct {
	encoded bitvector
	// blockSize is the number of bits in a block (value from 1-64)
	blockSize uint8
	// classSize (number of bits required to store the number of 1s for each
	// block, max: # of bits in the block)
	classSize uint8
	// offsetSize (number of bits required to store the offset for each block,
	// max: C(n, n/2) + 1)
	offsetSize uint8
}

type Integers interface {
	~int | ~uint |
		~int8 | ~uint8 |
		~int16 | ~uint16 |
		~int32 | ~uint32 |
		~int64 | ~uint64
}

// floor(log_2(n))
func floorLog2[T Integers](n T) (out uint8) {
	var zero T
	for n > zero {
		n >>= 1
		out++
	}
	return
}

func calcBlocksize(n uint64) (blocksize uint8) {
	blocksize = floorLog2(n)
	// blocksize / 2
	blocksize >>= 1
	return
}

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

func NewRRR(bits bitvector) (out RRR) {
	n := bits.Length()

	out.blockSize = calcBlocksize(n)
	out.classSize = floorLog2(out.blockSize) + 1

	maxOffset := choose(uint64(out.blockSize), uint64(out.blockSize)>>1)
	out.offsetSize = floorLog2(maxOffset) + 1

	blocks := n/uint64(out.blockSize) + 1
	totalSize := blocks * (uint64(out.classSize) + uint64(out.offsetSize))
	out.encoded = makeBitVector(totalSize)

	return
}
