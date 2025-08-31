package wavelettree

import "math/bits"

// RRR enables near practically O(1) calculations of bitwise rank(b, i) and
// select(b, i)
type RRR struct {
	bits BitVector
	// blockSize is the number of bits in a block (value from [1, 64])
	blockSize uint8
	// superblockSize stores the number of bits inside a super block.
	superblockSize uint16
	// classFieldSize is the number of bits required to store the number of 1s
	// for each block, max: # of bits in the block
	classFieldSize uint8
	// offsetFieldSize is the number of bits required to store the offset for
	// each block, max: C(n, n/2) + 1
	offsetFieldSize uint8
	// serializedBlockSize is the number of bits
	serializedBlockSize uint8
	// cumulativeRankFieldSize is the number of bits required to store the
	// cumulative rank of a superblock.
	cumulativeRankFieldSize uint8
	// serializedSuperblockSize is the number of bits a single super block takes up
	serializedSuperblockSize uint16
}

func rank[T uint8 | uint16 | uint32 | uint64](blocksize, bytesize, class uint8, content T) (offset uint64) {
	remaining := class
	mask := T(1)
	for range bytesize {
		if content&mask > 0 {
			offset += choose(uint64(blocksize-1), uint64(remaining))
		}
		mask <<= 1
	}
	return
}

func unrank[T uint8 | uint16 | uint32 | uint64]() {

}

func getBlockValues(blocksize uint8, i uint64, bitvec BitVector) (class uint8, offset uint64) {
	switch {
	case blocksize <= 8:
		content := bitvec.Get8(blocksize, i)
		class = uint8(bits.OnesCount8(content))
		offset = rank(blocksize, 8, class, content)
		return
	case blocksize <= 16:
		content := bitvec.Get16(blocksize, i)
		class = uint8(bits.OnesCount16(content))
		offset = rank(blocksize, 16, class, content)
		return
	case blocksize <= 32:
		content := bitvec.Get32(blocksize, i)
		class = uint8(bits.OnesCount32(content))
		offset = rank(blocksize, 32, class, content)
		return
	case blocksize <= 64:
		content := bitvec.Get64(blocksize, i)
		class = uint8(bits.OnesCount64(content))
		offset = rank(blocksize, 64, class, content)
		return
	}
	panic("exceeded max block length 64!")
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
	out.blockSize = blocksize

	superblocksize := opts.SuperBlockSize
	if opts.SuperBlockSize < 2 {
		// the max superblock size is 64
		superblocksize = nbitsize
	}
	// the max superblock size would be 64 * 64 = 4096
	out.superblockSize = uint16(superblocksize) * uint16(blocksize)

	// the size of the class field in the serialized block
	out.classFieldSize = floorLog2(blocksize)
	// the maximum possible value for offset (given by nCr(b, b/2))
	maxOffset := choose(uint64(blocksize), uint64(out.blockSize)>>1)
	// the size of the offset field in the serialized block
	out.offsetFieldSize = floorLog2(maxOffset)

	// worst case all ones up to the last superblock
	out.cumulativeRankFieldSize = floorLog2(n)

	blockNum := n / uint64(blocksize)
	// there is additional +1 because even if n cannot "fit" a single super
	// block, it will still be added at the start anyway
	superBlockNum := n/(uint64(blocksize)*uint64(superblocksize)) + 1

	// the serialized block size (in bits) of class + offset
	totalBlockSize := out.classFieldSize + out.offsetFieldSize
	out.serializedBlockSize = totalBlockSize
	// the total size (in bits) of the serialized block
	totalSize := blockNum*uint64(totalBlockSize) + superBlockNum*uint64(out.cumulativeRankFieldSize)
	out.bits = NewBitVector(totalSize)

	out.serializedSuperblockSize = uint16(out.cumulativeRankFieldSize) + uint16(totalBlockSize*superblocksize)

	// serialize blocks
	inCursor := uint64(0)
	outCursor := uint64(0)
	cumulativeRank := uint64(0)
	for i := range blockNum {
		if i%uint64(superblocksize) == 0 {
			switch {
			case out.cumulativeRankFieldSize <= 8:
				out.bits.Set8(out.cumulativeRankFieldSize, outCursor, uint8(cumulativeRank))
			case out.cumulativeRankFieldSize <= 16:
				out.bits.Set16(out.cumulativeRankFieldSize, outCursor, uint16(cumulativeRank))
			case out.cumulativeRankFieldSize <= 32:
				out.bits.Set32(out.cumulativeRankFieldSize, outCursor, uint32(cumulativeRank))
			case out.cumulativeRankFieldSize <= 64:
				out.bits.Set64(out.cumulativeRankFieldSize, outCursor, uint64(cumulativeRank))
			}
			outCursor += uint64(out.cumulativeRankFieldSize)
		}

		class, offset := getBlockValues(blocksize, inCursor, bits)
		inCursor += uint64(blocksize)

		cumulativeRank += uint64(class)

		// we know class field size will always be 2-3 bits
		out.bits.Set8(out.classFieldSize, outCursor, class)
		outCursor += uint64(out.classFieldSize)

		switch {
		case out.offsetFieldSize <= 8:
			out.bits.Set8(out.offsetFieldSize, outCursor, uint8(offset))
		case out.offsetFieldSize <= 16:
			out.bits.Set16(out.offsetFieldSize, outCursor, uint16(offset))
		case out.offsetFieldSize <= 32:
			out.bits.Set32(out.offsetFieldSize, outCursor, uint32(offset))
		case out.offsetFieldSize <= 64:
			out.bits.Set64(out.offsetFieldSize, outCursor, offset)
		}
		outCursor += uint64(out.offsetFieldSize)
	}

	return
}

// Rank returns the number of "bit" encountered from [0, i] in the bitvector
// where "bit" is either 0 or 1
func (r RRR) Rank(bit uint8, i uint64) uint64 {
	superblockIdx := i / uint64(r.superblockSize)
	superblockBitIdx := superblockIdx * uint64(r.serializedSuperblockSize)

	var rank uint64
	switch {
	case r.cumulativeRankFieldSize <= 8:
		rank = uint64(r.bits.Get8(r.cumulativeRankFieldSize, superblockBitIdx))
	case r.cumulativeRankFieldSize <= 16:
		rank = uint64(r.bits.Get16(r.cumulativeRankFieldSize, superblockBitIdx))
	case r.cumulativeRankFieldSize <= 32:
		rank = uint64(r.bits.Get32(r.cumulativeRankFieldSize, superblockBitIdx))
	case r.cumulativeRankFieldSize <= 64:
		rank = uint64(r.bits.Get64(r.cumulativeRankFieldSize, superblockBitIdx))
	}

	originalIdx := superblockIdx * uint64(r.superblockSize)
	cursor := superblockBitIdx + uint64(r.cumulativeRankFieldSize)
	for {
		rank += uint64(r.bits.Get8(r.classFieldSize, cursor))
		cursor += uint64(r.classFieldSize)
	}

	return rank
}
