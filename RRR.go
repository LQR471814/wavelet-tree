package wavelettree

// RRR enables near practically O(1) calculations of bitwise rank(b, i) and
// select(b, i)
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
	// superblockSize stores the number of blocks to include in a super block.
	superblockSize uint8
}

func computeOffset[T uint8 | uint16 | uint32 | uint64](blocksize, bytesize, class uint8, content T) (offset uint8) {
	remaining := class
	mask := T(1)
	for range bytesize {
		if content&mask > 0 {
			offset += uint8(choose(uint64(blocksize-1), uint64(remaining)))
		}
		mask <<= 1
	}
	return
}

func getBlockValues(blocksize uint8, i uint64, bits BitVector) (class, offset uint8) {
	switch {
	case blocksize <= 8:
		content := bits.Get8(blocksize, i)
		class = countbits(8, content)
		offset = computeOffset(blocksize, 8, class, content)
		return
	case blocksize <= 16:
		content := bits.Get16(blocksize, i)
		class = countbits(16, content)
		offset = computeOffset(blocksize, 16, class, content)
		return
	case blocksize <= 32:
		content := bits.Get32(blocksize, i)
		class = countbits(32, content)
		offset = computeOffset(blocksize, 32, class, content)
		return
	case blocksize <= 64:
		content := bits.Get64(blocksize, i)
		class = countbits(64, content)
		offset = computeOffset(blocksize, 64, class, content)
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

	superblocksize := opts.SuperBlockSize
	if opts.SuperBlockSize < 2 {
		superblocksize = nbitsize / blocksize
	}

	// the block size in bits to use on the input
	out.blockSize = blocksize

	// the size of the class field in the serialized block
	out.classFieldSize = floorLog2(blocksize)
	// the maximum possible value for offset (given by nCr(b, b/2))
	maxOffset := choose(uint64(blocksize), uint64(out.blockSize)>>1)
	// the size of the offset field in the serialized block
	out.offsetFieldSize = floorLog2(maxOffset)

	blockNum := n / uint64(blocksize)

	// the serialized block size (in bits) of class + offset
	serializedBlocksize := out.classFieldSize + out.offsetFieldSize
	// the total size (in bits) of the serialized block
	totalSerializedSize := blockNum * uint64(serializedBlocksize)
	out.encoded = NewBitVector(totalSerializedSize)

	// serialize blocks
	for i := range blockNum {
		bitIdx := i * uint64(blocksize)
		class, offset := getBlockValues(blocksize, bitIdx, bits)
		out.encoded.Set8(out.classFieldSize, bitIdx, class)
		out.encoded.Set8(out.offsetFieldSize, bitIdx+uint64(out.classFieldSize), offset)
	}

	return
}

// Rank returns the number of "bit" encountered from [0, i] in the bitvector
// where "bit" is either 0 or 1
func (r RRR) Rank(bit uint8, i uint64) uint64 {
	return 0
}
