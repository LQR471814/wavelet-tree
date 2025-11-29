package wavelettree

import (
	"fmt"
	"math/bits"
)

// RRR enables practically O(1) calculations of bitwise rank(b, i) and
// select(b, i)
type RRR struct {
	bits BitVector
	// blockLogicalSize is the logical number of bits in a block uncompressed
	// (value from [1, 64])
	blockLogicalSize uint8
	// superblockLogicalCapacity is the logical number of blocks contained
	// within a super block
	superblockLogicalCapacity uint8
	// superblockLogicalSize is the logical number of bits in a super block
	// uncompressed
	superblockLogicalSize uint16
	// classFieldActualSize is the number of bits required to store the class
	// field for each block. max: # of bits in the block (this one will always
	// be 2-3 bits)
	classFieldActualSize uint8
	// offsetFieldActualSize is the number of bits required to store the offset
	// for each block. max: C(n, n/2) + 1
	offsetFieldActualSize uint8
	// blockActualSize is the number of bits required to store a block in
	// actuality after compression
	blockActualSize uint8
	// cumulativeRankFieldActualSize is the number of bits required to store
	// the cumulative rank of a superblock.
	cumulativeRankFieldActualSize uint8
	// superblockActualSize is the number of bits required to store a
	// superblock in actuality after compression
	superblockActualSize uint16
}

// RRROptions allow you to configure some parameters of the RRR datastructure,
// usually you will not need to touch this
type RRROptions struct {
	// BlockSize defines the number of bits within a block
	//
	// It is a value from [1, 64], if 0 or unspecified it will automatically
	// calculate the theoretical optimal value and use it
	BlockSize uint8

	// SuperBlockSize defines the number of blocks within a super block
	//
	// It is a value from [2, 255], if < 2, it will automatically calculate the
	// theoretical optimal value and use it
	SuperBlockSize uint8
}

// NewRRR creates a new RRR datastructure
func NewRRR(bits BitVector, opts RRROptions) (out RRR) {
	n := bits.Length()
	nbitsize := floorLog2(n)

	blocksize := opts.BlockSize
	if opts.BlockSize > 64 {
		panic("blocksize must not be larger than 64!")
	}
	if opts.BlockSize == 0 {
		blocksize = nbitsize
		blocksize >>= 1
	}
	out.blockLogicalSize = blocksize

	out.superblockLogicalCapacity = opts.SuperBlockSize
	if opts.SuperBlockSize < 2 {
		// the max superblock size is 64
		out.superblockLogicalCapacity = nbitsize
	}
	// the max superblock size would be 64 * 64 = 4096
	out.superblockLogicalSize = uint16(out.superblockLogicalCapacity) * uint16(blocksize)

	// the size of the class field in the serialized block
	out.classFieldActualSize = floorLog2(blocksize)
	// the maximum possible value for offset (given by nCr(b, b/2))
	maxOffset := choose(uint64(blocksize), uint64(out.blockLogicalSize)>>1)
	// the size of the offset field in the serialized block
	out.offsetFieldActualSize = floorLog2(maxOffset)

	// worst case all ones up to the last superblock
	out.cumulativeRankFieldActualSize = floorLog2(n)

	blockNum := n / uint64(blocksize)
	// there is additional +1 because even if n cannot "fit" a single super
	// block, it will still be added at the start anyway
	superBlockNum := n/(uint64(blocksize)*uint64(out.superblockLogicalCapacity)) + 1

	// the serialized block size (in bits) of class + offset
	totalBlockSize := out.classFieldActualSize + out.offsetFieldActualSize
	out.blockActualSize = totalBlockSize
	// the total size (in bits) of the serialized block
	totalSize := blockNum*uint64(totalBlockSize) + superBlockNum*uint64(out.cumulativeRankFieldActualSize)
	out.bits = NewBitVector(totalSize)

	out.superblockActualSize = uint16(out.cumulativeRankFieldActualSize) + uint16(totalBlockSize*out.superblockLogicalCapacity)

	// serialize blocks
	inCursor := uint64(0)
	outCursor := uint64(0)
	cumulativeRank := uint64(0)
	for i := range blockNum {
		if i%uint64(out.superblockLogicalCapacity) == 0 {
			switch {
			case out.cumulativeRankFieldActualSize <= 8:
				out.bits.Set8(out.cumulativeRankFieldActualSize, outCursor, uint8(cumulativeRank))
			case out.cumulativeRankFieldActualSize <= 16:
				out.bits.Set16(out.cumulativeRankFieldActualSize, outCursor, uint16(cumulativeRank))
			case out.cumulativeRankFieldActualSize <= 32:
				out.bits.Set32(out.cumulativeRankFieldActualSize, outCursor, uint32(cumulativeRank))
			case out.cumulativeRankFieldActualSize <= 64:
				out.bits.Set64(out.cumulativeRankFieldActualSize, outCursor, uint64(cumulativeRank))
			}
			outCursor += uint64(out.cumulativeRankFieldActualSize)
		}

		class, offset := getBlockValues(blocksize, inCursor, bits)
		inCursor += uint64(blocksize)

		cumulativeRank += uint64(class)

		// we know class field size will always be 2-3 bits
		out.bits.Set8(out.classFieldActualSize, outCursor, class)
		outCursor += uint64(out.classFieldActualSize)

		switch {
		case out.offsetFieldActualSize <= 8:
			out.bits.Set8(out.offsetFieldActualSize, outCursor, uint8(offset))
		case out.offsetFieldActualSize <= 16:
			out.bits.Set16(out.offsetFieldActualSize, outCursor, uint16(offset))
		case out.offsetFieldActualSize <= 32:
			out.bits.Set32(out.offsetFieldActualSize, outCursor, uint32(offset))
		case out.offsetFieldActualSize <= 64:
			out.bits.Set64(out.offsetFieldActualSize, outCursor, offset)
		}
		outCursor += uint64(out.offsetFieldActualSize)
	}

	return
}

// Rank returns the number of 1-bits encountered from [0, i] in the bitvector.
//
// To get the number of 0-bits encountered, simply do i-rank(i).
func (r RRR) Rank(i uint64) uint64 {
	superblockIdx := i / uint64(r.superblockLogicalSize)
	remainder := i % uint64(r.superblockLogicalCapacity)
	blockIdx := remainder / uint64(r.blockLogicalSize)
	remainder = remainder % uint64(r.blockLogicalSize)

	cursor := superblockIdx * uint64(r.superblockActualSize)
	var ones uint64
	switch {
	case r.cumulativeRankFieldActualSize <= 8:
		ones = uint64(r.bits.Get8(r.cumulativeRankFieldActualSize, cursor))
	case r.cumulativeRankFieldActualSize <= 16:
		ones = uint64(r.bits.Get16(r.cumulativeRankFieldActualSize, cursor))
	case r.cumulativeRankFieldActualSize <= 32:
		ones = uint64(r.bits.Get32(r.cumulativeRankFieldActualSize, cursor))
	case r.cumulativeRankFieldActualSize <= 64:
		ones = uint64(r.bits.Get64(r.cumulativeRankFieldActualSize, cursor))
	}
	// ones = number of ones < the superblock containing "i"

	cursor += uint64(r.cumulativeRankFieldActualSize)
	idx := uint64(0)
	for {
		class := r.bits.Get8(r.classFieldActualSize, cursor)
		var offset uint64
		switch {
		case r.offsetFieldActualSize <= 8:
			offset = uint64(r.bits.Get8(r.offsetFieldActualSize, cursor+uint64(r.classFieldActualSize)))
		case r.offsetFieldActualSize <= 16:
			offset = uint64(r.bits.Get16(r.offsetFieldActualSize, cursor+uint64(r.classFieldActualSize)))
		case r.offsetFieldActualSize <= 32:
			offset = uint64(r.bits.Get32(r.offsetFieldActualSize, cursor+uint64(r.classFieldActualSize)))
		case r.offsetFieldActualSize <= 64:
			offset = uint64(r.bits.Get64(r.offsetFieldActualSize, cursor+uint64(r.classFieldActualSize)))
		}

		if idx == blockIdx {
			switch {
			case r.blockLogicalSize <= 8:
				block := unrank[uint8](class, offset)
				// mask away bits beyond the remainder after the target blockIdx
				block &= (^uint8(0)) >> (r.blockLogicalSize - uint8(remainder))
				ones += uint64(bits.OnesCount8(block))
			case r.blockLogicalSize <= 16:
				block := unrank[uint16](class, offset)
				block &= (^uint16(0)) >> (r.blockLogicalSize - uint8(remainder))
				ones += uint64(bits.OnesCount16(block))
			case r.blockLogicalSize <= 32:
				block := unrank[uint32](class, offset)
				block &= (^uint32(0)) >> (r.blockLogicalSize - uint8(remainder))
				ones += uint64(bits.OnesCount32(block))
			case r.blockLogicalSize <= 64:
				block := unrank[uint64](class, offset)
				block &= (^uint64(0)) >> (r.blockLogicalSize - uint8(remainder))
				ones += uint64(bits.OnesCount64(block))
			}
			break
		}
		ones += uint64(class)
		cursor += uint64(r.blockActualSize)
		idx++
	}

	return ones
}

// Select returns the i'th "bit" in the bitvector, where "bit" can either be 0
// or 1.
func (r RRR) Select(bit uint8, i uint64) {
	if bit != 0 && bit != 1 {
		panic(fmt.Errorf("invalid parameter for bit, must be either 0 or 1 got: %d", bit))
	}

	bitlength := r.bits.Length()
	lo := uint64(0)
	hi := bitlength / uint64(r.superblockActualSize) // index of last super block

	// find via binary search, the maximum superblock that
	for {
		mid := (lo + hi) / 2
	}

}

// rank computes the offset given a block of bits
func rank[T uint8 | uint16 | uint32 | uint64](blocksize uint8, content T) (offset uint64) {
	combIndex := 0
	mask := T(1)
	for pos := range blocksize {
		if content&mask > 0 {
			offset += choose(uint64(pos), uint64(combIndex+1))
			combIndex++
		}
		mask <<= 1
	}
	return
}

func unrank[T uint8 | uint16 | uint32 | uint64](class uint8, offset uint64) (content T) {
	// special case: 0
	if class == 0 {
		return 0
	}
	combIndex := class - 1
	remaining := offset
	for {
		var position uint8
		var lastcontrib uint64
		for pos := combIndex + 1; ; pos++ {
			contrib := choose(uint64(pos), uint64(combIndex+1))
			if contrib > uint64(remaining) {
				position = pos - 1
				remaining -= lastcontrib
				break
			}
			lastcontrib = contrib
		}
		content |= (1 << position)
		if combIndex == 0 {
			break
		}
		combIndex--
	}
	return
}

func getBlockValues(blocksize uint8, i uint64, bitvec BitVector) (class uint8, offset uint64) {
	switch {
	case blocksize <= 8:
		content := bitvec.Get8(blocksize, i)
		class = uint8(bits.OnesCount8(content))
		offset = rank(blocksize, content)
		return
	case blocksize <= 16:
		content := bitvec.Get16(blocksize, i)
		class = uint8(bits.OnesCount16(content))
		offset = rank(blocksize, content)
		return
	case blocksize <= 32:
		content := bitvec.Get32(blocksize, i)
		class = uint8(bits.OnesCount32(content))
		offset = rank(blocksize, content)
		return
	case blocksize <= 64:
		content := bitvec.Get64(blocksize, i)
		class = uint8(bits.OnesCount64(content))
		offset = rank(blocksize, content)
		return
	}
	panic("exceeded max block length 64!")
}
