package wavelettrie

import "fmt"

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
That is the result of `ceil(log_2(C(8, 4)))`.
*/

// bitvector is effectively a slice of bits
type bitvector []uint8

// Get gets the value of any size 1-8 bit word at any bit index i from the
// bitvector
func (v bitvector) Get(size uint8, i uint64) (result uint8) {
	if size < 1 || size > 8 {
		panic(fmt.Sprintf("invalid size: %d", size))
	}

	byte := i / 8
	bit := uint8(i % 8)

	var mask uint8 = 255 >> (8 - size)

	result = (v[byte] >> bit) & mask

	// 9 = 8+1 (since overlap threshold for size=1 should be 8)
	var overlapThreshold uint8 = 9 - size

	// overlap amount is: amount of bits to be found in the next byte after
	// bits are found in current byte
	var overlapAmount uint8 = size - (8 - bit)

	if bit >= overlapThreshold {
		next := v[byte+1]
		var nextmask uint8 = 255 >> overlapAmount
		var overlap uint8 = (next >> (overlapAmount - 1)) & nextmask
		result = result | overlap
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

	surrounding := v[byte] & (^mask)
	v[byte] = surrounding | (value << bit)

	// 9 = 8+1 (since overlap threshold for size=1 should be 8)
	var overlapThreshold uint8 = 9 - size

	// overlap amount is: amount of bits to be found in the next byte after
	// bits are found in current byte
	var bitsInCurrent uint8 = 8 - bit
	var overlapAmount uint8 = size - bitsInCurrent

	if bit >= overlapThreshold {
		next := v[byte+1]
		var nextMask uint8 = 255 >> overlapAmount
		nextupSurround := next & (^nextMask)
		// remove the bits in value that have already been set in the current
		// byte and set those bits in the next byte
		v[byte+1] = nextupSurround | (value >> bitsInCurrent)
	}
}

// 3 bits for class (represents 8 different possible 1s)
type block struct {
}

type RRR struct {
}
