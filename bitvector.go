package wavelettree

import (
	"fmt"
	"unsafe"
)

// BitVector is effectively a slice of bits.
type BitVector struct {
	bitlength uint64
	bytes     []byte
}

func NewBitVector(bitlength uint64) BitVector {
	bytelength := bitlength/8 + 1
	vec := BitVector{
		bitlength: bitlength,
		bytes:     make([]byte, bytelength),
	}
	return vec
}

// Get8 allows you to get 1-8 bits from the bitvector at once and return it as
// a uint8
func (v BitVector) Get8(size uint8, i uint64) uint8 {
	if size == 0 || size > 8 {
		panic(fmt.Sprintf("invalid bitsize: %d", size))
	}
	return getbits[uint8](8, size, v, i)
}

// Get16 allows you to get 1-16 bits from the bitvector at once and return it
// as a uint16
func (v BitVector) Get16(size uint8, i uint64) uint16 {
	if size == 0 || size > 16 {
		panic(fmt.Sprintf("invalid bitsize: %d", size))
	}
	return getbits[uint16](16, size, v, i)
}

// Get32 allows you to get 1-32 bits from the bitvector at once and return it
// as a uint32
func (v BitVector) Get32(size uint8, i uint64) uint32 {
	if size == 0 || size > 32 {
		panic(fmt.Sprintf("invalid bitsize: %d", size))
	}
	return getbits[uint32](32, size, v, i)
}

// Get64 allows you to get 1-64 bits from the bitvector at once and return it
// as a uint64
func (v BitVector) Get64(size uint8, i uint64) uint64 {
	if size == 0 || size > 64 {
		panic(fmt.Sprintf("invalid bitsize: %d", size))
	}
	return getbits[uint64](64, size, v, i)
}

// Set8 allows you to get 1-8 bits from the bitvector at once and return it as
// a uint8
func (v BitVector) Set8(size uint8, i uint64, value uint8) {
	if size == 0 || size > 8 {
		panic(fmt.Sprintf("invalid bitsize: %d", size))
	}
	setbits(8, size, v, i, value)
}

// Set16 allows you to get 1-16 bits from the bitvector at once and return it
// as a uint16
func (v BitVector) Set16(size uint8, i uint64, value uint16) {
	if size == 0 || size > 16 {
		panic(fmt.Sprintf("invalid bitsize: %d", size))
	}
	setbits(16, size, v, i, value)
}

// Set32 allows you to get 1-32 bits from the bitvector at once and return it
// as a uint32
func (v BitVector) Set32(size uint8, i uint64, value uint32) {
	if size == 0 || size > 32 {
		panic(fmt.Sprintf("invalid bitsize: %d", size))
	}
	setbits(32, size, v, i, value)
}

// Set64 allows you to get 1-64 bitsv.bytes from the bitvector at once and return it
// as a uint64
func (v BitVector) Set64(size uint8, i uint64, value uint64) {
	if size == 0 || size > 64 {
		panic(fmt.Sprintf("invalid bitsize: %d", size))
	}
	setbits(64, size, v, i, value)
}

// Append8 allows you to get 1-8 bits from the bitvector at once and return it as
// a uint8
func (v BitVector) Append8(size uint8, value uint8) BitVector {
	if size == 0 || size > 8 {
		panic(fmt.Sprintf("invalid bitsize: %d", size))
	}
	return appendbits(8, size, v, value)
}

// Append16 allows you to get 1-16 bits from the bitvector at once and return it
// as a uint16
func (v BitVector) Append16(size uint8, value uint16) BitVector {
	if size == 0 || size > 16 {
		panic(fmt.Sprintf("invalid bitsize: %d", size))
	}
	return appendbits(16, size, v, value)
}

// Append32 allows you to get 1-32 bits from the bitvector at once and return it
// as a uint32
func (v BitVector) Append32(size uint8, value uint32) BitVector {
	if size == 0 || size > 32 {
		panic(fmt.Sprintf("invalid bitsize: %d", size))
	}
	return appendbits(32, size, v, value)
}

// Append64 allows you to get 1-64 bits from the bitvector at once and return it
// as a uint64
func (v BitVector) Append64(size uint8, value uint64) BitVector {
	if size == 0 || size > 64 {
		panic(fmt.Sprintf("invalid bitsize: %d", size))
	}
	return appendbits(64, size, v, value)
}

// - bitsize can be any number of bits from 1-64
// - bytesize must be one of 8, 16, 32, or 64
func getbits[T uint8 | uint16 | uint32 | uint64](bytesize, bitsize uint8, v BitVector, i uint64) (result T) {
	if i > v.bitlength-uint64(bitsize) {
		panic(fmt.Sprintf("get bit index out of range: [%d]", i))
	}

	byteslice := *(*[]T)(unsafe.Pointer(&v.bytes))

	byte := i / uint64(bytesize)
	bitstart := uint8(i % uint64(bytesize))

	ALL_ONES := ^T(0)

	mask := ALL_ONES << bitstart
	retrieved := byteslice[byte]
	result = (retrieved & mask) >> bitstart

	// amount of bits set in the current byte
	var currentSet = bytesize - bitstart

	var overlapThreshold uint8 = (bytesize + 1) - bitsize
	var overlapAmount uint8 = bitsize - currentSet

	// ## params
	// - byte -> which byte to query from
	// - bitstart -> which bit in that byte to start setting
	// - overlapThreshold -> if bitstart > overlapThreshold then some bits will
	// spill over to the next byte
	// - overlapAmount -> amount of bits to be retrieved in the next byte

	// ## bit representations
	// - mask -> which bits in the current byte pertain to the bits to be retrieved
	// - retrieved -> bits in the current byte
	//
	// - nextMask -> which bits in the next byte pertain to bits that should be retrieved
	// - nextRetrieved -> bits in the next byte that have been retrieved

	// fmt.Printf(
	// 	"GET PARAMS | byte: %d | bitstart: %d | overlapThreshold: %d | overlapAmount: %d\n",
	// 	byte,
	// 	bitstart,
	// 	overlapThreshold,
	// 	overlapAmount,
	// )
	//
	// fmt.Printf(
	// 	fmt.Sprintf(
	// 		"GET BITS | mask: %%0%[1]db | retrieved: %%0%[1]db | final: %%0%[1]db\n",
	// 		bytesize,
	// 	),
	// 	mask,
	// 	retrieved,
	// 	result,
	// )

	if bitstart >= overlapThreshold {
		nextRetrieved := byteslice[byte+1]

		nextMask := ^(ALL_ONES << overlapAmount)
		overlap := nextRetrieved & nextMask
		result = result | (overlap << currentSet)

		// fmt.Printf(
		// 	fmt.Sprintf(
		// 		"GET NEXT BITS | mask: %%0%[1]db | nextRetrieved: %%0%[1]db | final: %%0%[1]db\n",
		// 		bytesize,
		// 	),
		// 	nextMask,
		// 	nextRetrieved,
		// 	result,
		// )
	}

	return
}

func setbits[T uint8 | uint16 | uint32 | uint64](bytesize, bitsize uint8, v BitVector, i uint64, value T) {
	// suppose bitlength = 25
	// bitsize = 8
	// the maximum index would be 25 - 8 = 17
	// that is, it would set indices [17, 24]
	if i > v.bitlength-uint64(bitsize) {
		panic(fmt.Sprintf("set bit index out of range: [%d]", i))
	}

	byteslice := *(*[]T)(unsafe.Pointer(&v.bytes))

	byte := i / uint64(bytesize)
	bitstart := uint8(i % uint64(bytesize))

	ALL_ONES := ^T(0)

	mask := ^(ALL_ONES << bitstart)
	newCurrent := value << bitstart

	original := byteslice[byte] & mask
	byteslice[byte] = original | newCurrent

	overlapThreshold := (bytesize + 1) - bitsize

	bitsSet := bytesize - bitstart
	overlapAmount := bitsize - bitsSet

	// ## params
	// - byte -> which byte to query from
	// - bitstart -> which bit in that byte to start setting bits from
	// - overlapThreshold -> if bitstart > overlapThreshold then some bits will
	// spill over to the next byte
	// - overlapAmount -> amount of bits to be set in the next byte

	// ## bit representations
	// - value -> the bits to set starting from `bitstart`
	//
	// - mask -> which bits in the current byte will NOT be set to the new value
	// - newCurrent -> the bits to update in the current byte
	//
	// - nextMask -> which bits in the next byte will NOT be set to the part of
	// the new value in the next byte
	// - nextNewCurrent -> the bits to update in the next byte

	// fmt.Printf(
	// 	"SET PARAMS | byte: %d | bitstart: %d | overlapThreshold: %d | overlapAmount: %d\n",
	// 	byte,
	// 	bitstart,
	// 	overlapThreshold,
	// 	overlapAmount,
	// )
	//
	// fmt.Printf(
	// 	fmt.Sprintf(
	// 		"SET BITS | value: %%0%[1]db | mask: %%0%[1]db | newCurrent: %%0%[1]db | final: %%0%[1]db\n",
	// 		bytesize,
	// 	),
	// 	value,
	// 	mask,
	// 	newCurrent,
	// 	original|newCurrent,
	// )

	if bitstart >= overlapThreshold {
		next := byteslice[byte+1]
		var nextMask T = ALL_ONES << overlapAmount
		nextupOriginal := next & nextMask
		nextNewCurrent := value >> bitsSet

		// fmt.Printf(
		// 	fmt.Sprintf(
		// 		"SET NEXT BITS | mask: %%0%[1]db | newCurrent: %%0%[1]db | final: %%0%[1]db\n",
		// 		bytesize,
		// 	),
		// 	nextMask,
		// 	nextNewCurrent,
		// 	nextupOriginal|nextNewCurrent,
		// )

		// remove the bits in value that have already been set in the current
		// byte and set those bits in the next byte
		byteslice[byte+1] = nextupOriginal | nextNewCurrent
	}
}

func appendbits[T uint8 | uint16 | uint32 | uint64](bytesize, bitsize uint8, v BitVector, value T) BitVector {
	byteslice := *(*[]T)(unsafe.Pointer(&v.bytes))
	v.bytes = *(*[]byte)(unsafe.Pointer(&byteslice))

	originalEnd := v.bitlength
	v.bitlength += uint64(bitsize)
	byteLen := v.bitlength/uint64(bytesize) + 1

	if int(byteLen) > len(byteslice) {
		byteslice = append(byteslice, 0)
	}
	setbits(bytesize, bitsize, v, originalEnd, value)

	return v
}

// Length returns the bit length of the bitvector.
func (v BitVector) Length() uint64 {
	return v.bitlength
}

func (v BitVector) String() string {
	out := make([]byte, 0, v.bitlength)
	for _, b := range v.bytes {
		mask := uint8(1)
		for range 8 {
			if mask&b > 0 {
				out = append(out, '1')
			} else {
				out = append(out, '0')
			}
			mask <<= 1
		}
	}
	return string(out[:v.bitlength])
}
